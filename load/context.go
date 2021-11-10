package load

import (
	"context"
	"harkonnen/runtime"
)

type Context struct {
	context.Context
	*runtime.VariablePool
	*runtime.ErrorHandler
	*runtime.SampleCollector
	cancelFunc context.CancelFunc
}

func NewContext(parent context.Context, variablePool *runtime.VariablePool, errorHandler *runtime.ErrorHandler, sampleCollector *runtime.SampleCollector) *Context {
	output := new(Context)
	output.VariablePool = variablePool
	output.ErrorHandler = errorHandler
	output.SampleCollector = sampleCollector

	output.Context, output.cancelFunc = context.WithCancel(parent)

	return output
}

func (c *Context) Cancel() {
	c.cancelFunc()
}
