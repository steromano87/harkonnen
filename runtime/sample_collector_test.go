package runtime

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockedSample struct {
	mock.Mock
}

func (mocked MockedSample) Name() string {
	return "Mocked sample"
}

func (mocked MockedSample) Start() time.Time {
	return time.Time{}
}

func (mocked MockedSample) End() time.Time {
	return time.Now()
}

func (mocked MockedSample) Duration() time.Duration {
	return time.Now().Sub(time.Time{})
}

func (mocked MockedSample) SentBytes() int64 {
	return int64(256)
}

func (mocked MockedSample) ReceivedBytes() int64 {
	return int64(2048)
}

func TestNewSampleCollector(t *testing.T) {
	collector := NewSampleCollector()
	assert.IsType(t, &SampleCollector{}, collector)

	flushedSamples := collector.Flush()
	assert.Empty(t, flushedSamples, "Flushing unused sample collector should give no collected samples")
}

func TestSampleCollector_CollectFlush(t *testing.T) {
	collector := NewSampleCollector()
	mockedSample := MockedSample{}
	collector.Collect(mockedSample)

	flushedSamples := collector.Flush()
	if assert.Equal(t, 1, len(flushedSamples), "Expected only one sample") {
		assert.Equal(t, mockedSample, flushedSamples[0])
	}
}
