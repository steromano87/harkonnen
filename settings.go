package main

import "harkonnen/rest"

type Settings struct {
	Http rest.Settings `mapstructure:"http"`
}
