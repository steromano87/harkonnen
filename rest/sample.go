package rest

import (
	"github.com/steromano87/harkonnen/telemetry"
	"net/url"
	"time"
)

type Sample struct {
	telemetry.BaseSample
	URL        *url.URL
	Parameters url.Values
	Method     string
	IsRedirect bool
	FinalURL   *url.URL
}

func NewSample(name string, start time.Time, end time.Time, sentBytes int64, receivedBytes int64) Sample {
	sample := new(Sample)
	sample.BaseSample = telemetry.NewBaseSample(name, start, end, sentBytes, receivedBytes)

	return *sample
}
