package compiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

type Compiler struct {
	GoExecPath                string
	GoVersion                 string
	TempBuildCachePath        string
	CompiledScriptsFolderPath string
	CompiledScripts           []*File
}

const (
	BuildCachePrefix                 = "harkonnen-build-cache-"
	DefaultCompiledScriptsFolderPath = "./runtime/build"
	FakeMainFileName                 = "main.go"
)

func NewCompiler(goExecPath string, compiledScriptsFolderPath string) (*Compiler, error) {
	compiler := new(Compiler)
	var err error

	// Check Go executable
	err = compiler.checkGoExecutable(goExecPath)
	if err != nil {
		return nil, err
	}

	// Check Go version
	err = compiler.checkGoVersion()
	if err != nil {
		return nil, err
	}

	if compiledScriptsFolderPath != "" {
		compiler.CompiledScriptsFolderPath = compiledScriptsFolderPath
	} else {
		compiler.CompiledScriptsFolderPath = DefaultCompiledScriptsFolderPath
	}

	return compiler, nil
}

func (compiler *Compiler) Compile(script *File) error {
	var err error

	// Create temporary build cache (to be deleted right after the build is done)
	err = compiler.createBuildCache()
	defer func() {
		_ = compiler.cleanBuildCache()
	}()
	if err != nil {
		return err
	}

	// Read file content
	fileContent, err := ioutil.ReadFile(script.Path)
	if err != nil {
		return err
	}

	// Create fake main.go file
	err = compiler.dumpScriptInFakeMainFile(fileContent)
	if err != nil {
		return err
	}

	// Compile file as plugin
	compiledScriptPath, err := compiler.doCompile(script)
	if err != nil {
		return err
	}

	// Move compiled plugin to compiled scripts folder
	finalPath, err := compiler.moveCompiledScriptToCompiledScriptsFolder(compiledScriptPath)
	if err != nil {
		return err
	}

	script.CompiledObjectPath = finalPath
	compiler.CompiledScripts = append(compiler.CompiledScripts, script)
	return nil
}

func (compiler *Compiler) Clean() error {
	compiler.CompiledScripts = []*File{}
	return os.RemoveAll(compiler.CompiledScriptsFolderPath)
}

func (compiler *Compiler) checkGoExecutable(manuallySetGoExecutable string) error {
	// If path is manually defined, check if it exists
	if manuallySetGoExecutable != "" {
		if _, err := os.Stat(compiler.GoExecPath); os.IsNotExist(err) {
			return ErrGoCompilerNotFound{ManuallySetPath: compiler.GoExecPath}
		}

		compiler.GoExecPath = manuallySetGoExecutable
		return nil
	}

	// If no path is provided, detect where Go executable us
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("where go")
	case "linux", "freebsd", "darwin":
		cmd = exec.Command("which go")
	default:
		return ErrUnsupportedOS{OSName: runtime.GOOS}
	}

	output, err := compiler.executeCrossPlatformCommand(cmd)
	if err != nil {
		return ErrGoCompilerNotFound{}
	}

	compiler.GoExecPath = string(output)
	return nil
}

func (compiler *Compiler) checkGoVersion() error {
	// Get the whole version string
	compiler.GoVersion = runtime.Version()

	// TODO: add check for minimum Go version
	return nil
}

func (compiler *Compiler) createBuildCache() error {
	var err error
	compiler.TempBuildCachePath, err = ioutil.TempDir("", BuildCachePrefix)
	return err
}

func (compiler *Compiler) cleanBuildCache() error {
	return os.RemoveAll(compiler.TempBuildCachePath)
}

func (compiler *Compiler) dumpScriptInFakeMainFile(scriptContent []byte) error {
	fakeMainPath := path.Join(compiler.TempBuildCachePath, FakeMainFileName)
	return ioutil.WriteFile(fakeMainPath, scriptContent, 0644)
}

func (compiler *Compiler) doCompile(script *File) (string, error) {
	pluginExtension, err := compiler.pluginExtension()

	if err != nil {
		return "", err
	}

	cmd := exec.Command(
		fmt.Sprintf(
			"cd %s && %s build -buildmode=plugin -o %s%s %s",
			compiler.TempBuildCachePath,
			compiler.GoExecPath,
			script.Hash,
			pluginExtension,
			FakeMainFileName))

	output, err := compiler.executeCrossPlatformCommand(cmd)
	if err != nil {
		return "", ErrCompilationFailure{Message: string(output)}
	}

	return path.Join(compiler.TempBuildCachePath, script.Hash+pluginExtension), nil
}

func (compiler *Compiler) moveCompiledScriptToCompiledScriptsFolder(compiledScriptPath string) (string, error) {
	finalCompiledScriptPath := path.Join(compiler.CompiledScriptsFolderPath, path.Base(compiledScriptPath))
	err := os.Rename(compiledScriptPath, finalCompiledScriptPath)

	return finalCompiledScriptPath, err
}

func (compiler *Compiler) executeCrossPlatformCommand(command *exec.Cmd) ([]byte, error) {
	var execCmd *exec.Cmd
	commandString := strings.ReplaceAll(command.String(), "\n", "")

	switch runtime.GOOS {
	case "windows":
		execCmd = exec.Command("cmd", "/C", commandString)
	case "linux", "freebsd", "darwin":
		execCmd = exec.Command("sh", "-c", commandString)
	default:
		return []byte{}, ErrUnsupportedOS{OSName: runtime.GOOS}
	}

	return execCmd.CombinedOutput()
}

func (compiler Compiler) pluginExtension() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return ".dll", nil
	case "linux", "freebsd":
		return ".so", nil
	case "darwin":
		return ".dylib", nil
	default:
		return "", ErrUnsupportedOS{OSName: runtime.GOOS}
	}
}
