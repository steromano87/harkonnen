package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVariableNotFound_Error(t *testing.T) {
	myError := VariableNotFound{Name: "variableName"}

	assert.EqualError(t, myError, "variable 'variableName' not found", "Wrong error message format")
}
