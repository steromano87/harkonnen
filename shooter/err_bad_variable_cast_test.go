package shooter_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/shooter"
	"testing"
)

func TestErrBadVariableCast_Error(t *testing.T) {
	myError := &shooter.ErrBadVariableCast{
		Name:     "variableName",
		CastType: "int",
		RawValue: "true",
	}

	assert.EqualError(
		t,
		myError,
		"error when casting variable 'variableName' as int, raw value is 'true'",
		"Wrong error message format")
}
