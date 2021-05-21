package errors

import (
	"fmt"
)

type VariableCast struct {
	Name     string
	CastType string
	RawValue interface{}
}

func (vc VariableCast) Error() string {
	return fmt.Sprintf("error when casting variable '%s' as %s, raw value is '%v'", vc.Name, vc.CastType, vc.RawValue)
}
