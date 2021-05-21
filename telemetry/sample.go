package telemetry

import "time"

type Sample interface {
	Name() string
	Start() time.Time
	End() time.Time
	Duration() time.Duration
	SentBytes() int64
	ReceivedBytes() int64
}
