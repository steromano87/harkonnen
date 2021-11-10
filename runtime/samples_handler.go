package runtime

type SamplesHandler interface {
	Collect(sample Sample)
	Flush() []Sample
}
