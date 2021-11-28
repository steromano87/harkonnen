package log

type Collector struct {
	logs []Entry
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
