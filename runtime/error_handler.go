package runtime

type ErrorHandler struct {
	collectedErrors []error
}

func (handler *ErrorHandler) Capture(err error) {
	if err != nil {
		handler.collectedErrors = append(handler.collectedErrors, err)
	}
}

func (handler *ErrorHandler) GetCollected() []error {
	return handler.collectedErrors
}
