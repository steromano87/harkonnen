package compiler

import "fmt"

type ErrCompilationFailure struct {
	Message string
}

func (ce ErrCompilationFailure) Error() string {
	return fmt.Sprintf("got an error during compilation: '%s'", ce.Message)
}
