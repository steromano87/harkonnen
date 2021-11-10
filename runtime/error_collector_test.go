package runtime_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"harkonnen/runtime"
	"testing"
)

func TestErrorCollector_Capture_GetCollected_Single(t *testing.T) {
	collector := runtime.ErrorCollector{}
	err := fmt.Errorf("sample error: %v", "example")
	collector.Capture(err)
	collectedErrors := collector.GetCollected()

	if assert.Equal(t, 1, len(collectedErrors)) {
		assert.Equal(t, err, collectedErrors[0], "Captured error and collected error mismatch")
	}
}

func TestErrorCollector_Capture_GetCollected_Nil(t *testing.T) {
	collector := runtime.ErrorCollector{}
	collector.Capture(nil)
	collectedErrors := collector.GetCollected()

	assert.Empty(t, collectedErrors, "nil errors should not be collected")
}

func TestErrorCollector_Capture_GetCollected_Multiple(t *testing.T) {
	collector := runtime.ErrorCollector{}
	err1 := fmt.Errorf("first sample error: %v", "example")
	collector.Capture(err1)
	err2 := fmt.Errorf("second sample error: %v", "example")
	collector.Capture(err2)

	collectedErrors := collector.GetCollected()

	if assert.Equal(t, 2, len(collectedErrors), "Expected 2 errors to be collected") {
		assert.Equal(t, err1, collectedErrors[0])
		assert.Equal(t, err2, collectedErrors[1])
	}
}

func TestErrorCollector_HasErrors(t *testing.T) {
	collector := runtime.ErrorCollector{}
	err := fmt.Errorf("sample error")
	collector.Capture(err)

	assert.True(t, collector.HasErrors())
}
