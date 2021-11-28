package log_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/log"
	"testing"
	"time"
)

func TestCollector_Collect_Flush(t *testing.T) {
	collector := log.Collector{}

	collector.Collect(log.Entry{
		Timestamp: time.Time{},
		Level:     log.Error,
		Message:   "An error",
	})

	collector.Collect(log.Entry{
		Timestamp: time.Time{},
		Level:     log.Warning,
		Message:   "A warning",
	})

	collector.Collect(log.Entry{
		Timestamp: time.Time{},
		Level:     log.Info,
		Message:   "An info",
	})

	collector.Collect(log.Entry{
		Timestamp: time.Time{},
		Level:     log.Debug,
		Message:   "A debug message",
	})

	assert.Equal(t, 1, len(collector.Flush(log.Error)))
	assert.Equal(t, 2, len(collector.Flush(log.Warning)))
	assert.Equal(t, 3, len(collector.Flush(log.Info)))
	assert.Equal(t, 4, len(collector.Flush(log.Debug)))
}
