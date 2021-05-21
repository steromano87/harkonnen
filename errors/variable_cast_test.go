package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVariableCast_Error(t *testing.T) {
	myError := &VariableCast{
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
