package log_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/log"
	"testing"
)

func TestLevel_Above(t *testing.T) {
	assert.True(t, log.Error.Above(log.Warning))
	assert.False(t, log.Debug.Above(log.Info))
}
