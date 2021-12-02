package compiler

import "fmt"

type ErrUnsupportedOS struct {
	OSName string
}

func (uo ErrUnsupportedOS) Error() string {
	return fmt.Sprintf(
		"'%s' is not a supported OS. Supported OSes are: Windows, Linux, FreeBSD, Darwin (macOS)",
		uo.OSName)
}
