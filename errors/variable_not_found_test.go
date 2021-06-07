package errors_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/errors"
	"testing"
)

func TestVariableNotFound_Error(t *testing.T) {
	myError := errors.VariableNotFound{Name: "variableName"}

	assert.EqualError(t, myError, "variable 'variableName' not found", "Wrong error message format")
}
