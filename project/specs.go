package project

import (
	"github.com/steromano87/harkonnen/load"
)

type Specs struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`

	Scripts struct {
		SetUp    ScriptFile   `yaml:"set_up,omitempty"`
		Main     []ScriptFile `yaml:"main"`
		TearDown ScriptFile   `yaml:"tear_down,omitempty"`
	} `yaml:"scripts"`

	Ramps []load.LinearRamp `yaml:"ramps"`
}

func (s Specs) File() string {
	return "Harkonnen.yaml"
}

func (s Specs) Create(workdir string) error {
	return nil
}
