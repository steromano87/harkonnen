package runtime

import (
	"context"
)

type Context struct {
	context.Context
	*VariablePool
	*ErrorCollector
	*SampleCollector
	cancelFunc context.CancelFunc
}

func NewContext(parent context.Context, variablePool *VariablePool, errorHandler *ErrorCollector, sampleCollector *SampleCollector) Context {
	output := new(Context)
	output.VariablePool = variablePool
	output.ErrorCollector = errorHandler
	output.SampleCollector = sampleCollector

	output.Context, output.cancelFunc = context.WithCancel(parent)

	return *output
}

func (c *Context) Cancel() {
	c.cancelFunc()
}

func (c Context) NextLoop() <-chan struct{} {
	return make(chan struct{})
}
