package errors

import "fmt"

type VariableNotFound struct {
	Name string
}

func (vnf VariableNotFound) Error() string {
	return fmt.Sprintf("variable '%s' not found", vnf.Name)
}
