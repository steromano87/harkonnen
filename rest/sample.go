package rest

import (
	"harkonnen/telemetry"
	"net/url"
	"time"
)

type Sample struct {
	telemetry.Sample
	URL        *url.URL
	Parameters url.Values
	Method     string
	IsRedirect bool
	FinalURL   *url.URL
}

func NewSample(name string, start time.Time, end time.Time, sentBytes int64, receivedBytes int64) Sample {
	sample := new(Sample)
	sample.Sample = telemetry.NewSample(name, start, end, sentBytes, receivedBytes)

	return *sample
}
