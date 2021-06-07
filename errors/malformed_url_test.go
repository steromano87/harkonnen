package errors_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/errors"
	"testing"
)

func TestMalformedUrl_Error(t *testing.T) {
	myError := errors.MalformedUrl{Url: "my URL"}

	assert.EqualError(t, myError, "'my URL' is a malformed URL")
}
