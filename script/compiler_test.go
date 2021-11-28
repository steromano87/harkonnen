package script_test

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"harkonnen/script"
	"os"
	"testing"
)

func TestNewCompiler_AutoDetect(t *testing.T) {
	compiler, err := script.NewCompiler("", "")
	if assert.NoError(t, err) {
		assert.IsType(t, &script.Compiler{}, compiler)

		assert.NotEmpty(t, compiler.GoExecPath, "got empty Go executable path")
		assert.Regexp(t, "\\/.+", compiler.GoExecPath)
		assert.NotEmpty(t, compiler.GoVersion, "got empty Go version")
		assert.Regexp(t, "\\d\\.\\d+", compiler.GoVersion)
	}
}

func TestCompiler_Compile(t *testing.T) {
	fileContent := `
	package main

	func Test() string {
		return "Hello world"
	}
`

	compiledScriptsFolder := filet.TmpDir(t, os.TempDir())
	testFile := filet.TmpFile(t, os.TempDir(), fileContent)
	defer filet.CleanUp(t)

	compiler, err := script.NewCompiler("", compiledScriptsFolder)

	if assert.NoError(t, err) {
		testScript, _ := script.NewCompilableFile(testFile.Name(), "test")

		err := compiler.Compile(testScript)
		if assert.NoError(t, err) {
			assert.FileExists(t, testScript.CompiledObjectPath)
			assert.Len(t, compiler.CompiledScripts, 1)
			assert.NoFileExists(t, compiler.TempBuildCachePath)
		}
	}
}
