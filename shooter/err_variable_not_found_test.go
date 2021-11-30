package shooter_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/shooter"
	"testing"
)

func TestErrVariableNotFound_Error(t *testing.T) {
	myError := shooter.ErrVariableNotFound{Name: "variableName"}

	assert.EqualError(t, myError, "variable 'variableName' not found", "Wrong error message format")
}
