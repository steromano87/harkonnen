package runtime

type SampleCollector struct {
	samples []Sample
}

func NewSampleCollector() *SampleCollector {
	collector := new(SampleCollector)
	collector.samples = []Sample{}
	return collector
}

func (collector *SampleCollector) Collect(sample Sample) {
	collector.samples = append(collector.samples, sample)
}

func (collector *SampleCollector) Flush() []Sample {
	output := collector.samples
	collector.samples = []Sample{}
	return output
}
