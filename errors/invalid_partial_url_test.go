package errors_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/errors"
	"testing"
)

func TestInvalidPartialUrl_Error(t *testing.T) {
	myError := errors.InvalidPartialUrl{Url: "invalid URL"}

	assert.EqualError(
		t,
		myError,
		"'invalid URL' is an invalid partial URL",
		"Wrong error message format")
}
