package shooter_test

import (
	"github.com/steromano87/harkonnen/shooter"
	"github.com/stretchr/testify/assert"
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
