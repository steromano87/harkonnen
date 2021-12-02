package compiler_test

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"harkonnen/compiler"
	"os"
	"testing"
)

func TestNewCompiler_AutoDetect(t *testing.T) {
	comp, err := compiler.NewCompiler("", "")
	if assert.NoError(t, err) {
		assert.IsType(t, &compiler.Compiler{}, comp)

		assert.NotEmpty(t, comp.GoExecPath, "got empty Go executable path")
		assert.Regexp(t, "\\/.+", comp.GoExecPath)
		assert.NotEmpty(t, comp.GoVersion, "got empty Go version")
		assert.Regexp(t, "\\d\\.\\d+", comp.GoVersion)
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

	comp, err := compiler.NewCompiler("", compiledScriptsFolder)

	if assert.NoError(t, err) {
		testScript, _ := compiler.NewCompilableFile(testFile.Name(), "test")

		err := comp.Compile(testScript)
		if assert.NoError(t, err) {
			assert.FileExists(t, testScript.CompiledObjectPath)
			assert.Len(t, comp.CompiledScripts, 1)
			assert.NoFileExists(t, comp.TempBuildCachePath)
		}
	}
}
