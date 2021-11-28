package rest_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/rest"
	"testing"
)

func TestInvalidPartialUrl_Error(t *testing.T) {
	myError := rest.InvalidPartialUrlError{Url: "invalid URL"}

	assert.EqualError(
		t,
		myError,
		"'invalid URL' is an invalid partial URL",
		"Wrong error message format")
}
