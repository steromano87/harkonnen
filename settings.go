package harkonnen

import "harkonnen/httpp"

type Settings struct {
	Http httpp.Settings `mapstructure:"http"`
}
