package compiler

import "fmt"

type ErrGoCompilerNotFound struct {
	ManuallySetPath string
}

func (gcnf ErrGoCompilerNotFound) Error() string {
	if gcnf.ManuallySetPath != "" {
		return fmt.Sprintf("'%s' is not a valid Go executable path", gcnf.ManuallySetPath)
	} else {
		return "cannot find Go executable in PATH"
	}
}
