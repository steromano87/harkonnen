package main

import "github.com/steromano87/harkonnen/rest"

type Settings struct {
	Http rest.Settings `mapstructure:"http"`
}
