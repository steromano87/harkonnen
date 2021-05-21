package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorHandler_Capture_GetCollected_Single(t *testing.T) {
	handler := ErrorHandler{}
	err := fmt.Errorf("sample error: %v", "example")
	handler.Capture(err)
	collectedErrors := handler.GetCollected()

	if assert.Equal(t, 1, len(collectedErrors)) {
		assert.Equal(t, err, collectedErrors[0], "Captured error and collected error mismatch")
	}
}

func TestErrorHandler_Capture_GetCollected_Nil(t *testing.T) {
	handler := ErrorHandler{}
	handler.Capture(nil)
	collectedErrors := handler.GetCollected()

	assert.Empty(t, collectedErrors, "nil errors should not be collected")
}

func TestErrorHandler_Capture_GetCollected_Multiple(t *testing.T) {
	handler := ErrorHandler{}
	err1 := fmt.Errorf("first sample error: %v", "example")
	handler.Capture(err1)
	err2 := fmt.Errorf("second sample error: %v", "example")
	handler.Capture(err2)

	collectedErrors := handler.GetCollected()

	if assert.Equal(t, 2, len(collectedErrors), "Expected 2 errors to be collected") {
		assert.Equal(t, err1, collectedErrors[0])
		assert.Equal(t, err2, collectedErrors[1])
	}
}
