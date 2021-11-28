package shooter_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/shooter"
	"testing"
)

func TestVariableNotFound_Error(t *testing.T) {
	myError := shooter.VariableNotFoundError{Name: "variableName"}

	assert.EqualError(t, myError, "variable 'variableName' not found", "Wrong error message format")
}
