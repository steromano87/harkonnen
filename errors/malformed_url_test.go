package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMalformedUrl_Error(t *testing.T) {
	myError := MalformedUrl{Url: "my URL"}

	assert.EqualError(t, myError, "'my URL' is a malformed URL")
}
