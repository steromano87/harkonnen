package runtime

type ErrorCollector struct {
	collectedErrors []error
}

func (handler *ErrorCollector) Capture(err error) {
	if err != nil {
		handler.collectedErrors = append(handler.collectedErrors, err)
	}
}

func (handler *ErrorCollector) GetCollected() []error {
	return handler.collectedErrors
}

func (handler ErrorCollector) HasErrors() bool {
	return len(handler.collectedErrors) > 0
}
