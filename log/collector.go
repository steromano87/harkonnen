package log

import "time"

type Collector struct {
	logs []Entry
}

func (c *Collector) Error(message string) {
	c.Collect(Entry{
		Timestamp: time.Time{},
		Level:     ErrorLevel,
		Message:   message,
	})
}

func (c *Collector) Warning(message string) {
	c.Collect(Entry{
		Timestamp: time.Time{},
		Level:     WarningLevel,
		Message:   message,
	})
}

func (c *Collector) Info(message string) {
	c.Collect(Entry{
		Timestamp: time.Time{},
		Level:     InfoLevel,
		Message:   message,
	})
}

func (c *Collector) Debug(message string) {
	c.Collect(Entry{
		Timestamp: time.Time{},
		Level:     DebugLevel,
		Message:   message,
	})
}

func (c *Collector) Collect(entry Entry) {
	c.logs = append(c.logs, entry)
}

func (c *Collector) Flush(level Level) []Entry {
	var output []Entry

	for _, entry := range c.logs {
		if entry.Level.Above(level) {
			output = append(output, entry)
		}
	}

	return output
}
