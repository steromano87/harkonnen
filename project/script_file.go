package project

import (
	"errors"
	"os"
)

type ScriptFile struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func (f ScriptFile) Exists() bool {
	_, err := os.Stat(f.Path)
	return !errors.Is(err, os.ErrNotExist)
}
