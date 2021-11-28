package log_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/log"
	"testing"
)

func TestLevel_Above(t *testing.T) {
	assert.True(t, log.ErrorLevel.Above(log.WarningLevel))
	assert.False(t, log.DebugLevel.Above(log.InfoLevel))
}
