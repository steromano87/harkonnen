package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidPartialUrl_Error(t *testing.T) {
	myError := InvalidPartialUrl{Url: "invalid URL"}

	assert.EqualError(
		t,
		myError,
		"'invalid URL' is an invalid partial URL",
		"Wrong error message format")
}
