package shooter_test

import (
	"github.com/steromano87/harkonnen/shooter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrVariableNotFound_Error(t *testing.T) {
	myError := shooter.ErrVariableNotFound{Name: "variableName"}

	assert.EqualError(t, myError, "variable 'variableName' not found", "Wrong error message format")
}
