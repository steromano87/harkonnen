package telemetry

type SampleCollector struct {
	samples []Sample
}

func (collector *SampleCollector) Collect(sample Sample) {
	collector.samples = append(collector.samples, sample)
}

func (collector *SampleCollector) Flush() []Sample {
	output := collector.samples
	collector.samples = []Sample{}
	return output
}
