package log_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/log"
	"testing"
)

func TestCollector_Collect_Flush(t *testing.T) {
	collector := log.Collector{}

	collector.Error("An error")
	collector.Warning("A warning")
	collector.Info("An info")
	collector.Debug("A debug message")

	assert.Equal(t, 1, len(collector.Flush(log.ErrorLevel)))
	assert.Equal(t, 2, len(collector.Flush(log.WarningLevel)))
	assert.Equal(t, 3, len(collector.Flush(log.InfoLevel)))
	assert.Equal(t, 4, len(collector.Flush(log.DebugLevel)))
}
