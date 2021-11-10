package runtime

type ErrorsHandler interface {
	Capture(err error)
	GetCollected() []error
	HasErrors() bool
}
