package http

import (
	"time"
)

type Sample struct {
	start         time.Time
	end           time.Time
	name          string
	sentBytes     int64
	receivedBytes int64
	Info          SampleInfo
}

func (sample Sample) Name() string {
	return sample.name
}

func (sample Sample) Start() time.Time {
	return sample.start
}

func (sample Sample) End() time.Time {
	return sample.end
}

func (sample Sample) Duration() time.Duration {
	return sample.end.Sub(sample.start)
}

func (sample Sample) SentBytes() int64 {
	return sample.sentBytes
}

func (sample Sample) ReceivedBytes() int64 {
	return sample.receivedBytes
}
