package telemetry

import "time"

type BaseSample struct {
	start         time.Time
	end           time.Time
	name          string
	sentBytes     int64
	receivedBytes int64
}

func NewBaseSample(name string, start time.Time, end time.Time, sentBytes int64, receivedBytes int64) BaseSample {
	sample := new(BaseSample)
	sample.name = name
	sample.start = start
	sample.end = end
	sample.sentBytes = sentBytes
	sample.receivedBytes = receivedBytes

	return *sample
}

func (sample BaseSample) Name() string {
	return sample.name
}

func (sample BaseSample) Start() time.Time {
	return sample.start
}

func (sample BaseSample) End() time.Time {
	return sample.end
}

func (sample BaseSample) Duration() time.Duration {
	return sample.end.Sub(sample.start)
}

func (sample BaseSample) SentBytes() int64 {
	return sample.sentBytes
}

func (sample BaseSample) ReceivedBytes() int64 {
	return sample.receivedBytes
}
