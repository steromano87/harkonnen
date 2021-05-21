package runtime

import (
	"harkonnen/telemetry"
)

type SampleCollector struct {
	samples []telemetry.Sample
}

func NewSampleCollector() *SampleCollector {
	collector := new(SampleCollector)
	collector.samples = []telemetry.Sample{}
	return collector
}

func (collector *SampleCollector) Collect(sample telemetry.Sample) {
	collector.samples = append(collector.samples, sample)
}

func (collector *SampleCollector) Flush() []telemetry.Sample {
	output := collector.samples
	collector.samples = []telemetry.Sample{}
	return output
}
