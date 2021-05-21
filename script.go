package harkonnen

import "harkonnen/runtime"

type Script interface {
	Run(runtime *runtime.Runtime) error
}
