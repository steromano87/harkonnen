package shooter

import (
	"context"
	"harkonnen/log"
	"harkonnen/telemetry"
)

type Context struct {
	context.Context

	variablePool    *VariablePool
	sampleCollector *telemetry.SampleCollector
	logCollector    *log.Collector

	cancelFunc context.CancelFunc
}

func NewContext(parent context.Context) Context {
	output := new(Context)
	output.sampleCollector = new(telemetry.SampleCollector)
	output.logCollector = new(log.Collector)
	output.variablePool = NewVariablePool(output.logCollector)

	output.Context, output.cancelFunc = context.WithCancel(parent)

	return *output
}

func (c *Context) VariablePool() *VariablePool {
	return c.variablePool
}

func (c *Context) SampleCollector() *telemetry.SampleCollector {
	return c.sampleCollector
}

func (c *Context) LogCollector() *log.Collector {
	return c.logCollector
}

func (c *Context) Cancel() {
	c.cancelFunc()
}

func (c *Context) NextLoop() <-chan struct{} {
	return make(chan struct{})
}

func (c *Context) OnUnrecoverableError(err error) {
	c.logCollector.Error(err.Error())
	panic(err)
}
