package errors_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/errors"
	"testing"
)

func TestVariableCast_Error(t *testing.T) {
	myError := &errors.VariableCast{
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
