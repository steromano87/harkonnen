package harkonnen

import "harkonnen/http"

type Settings struct {
	Http http.Settings `mapstructure:"http"`
}
