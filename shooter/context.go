package shooter

import (
	"context"
	"github.com/rs/zerolog"
	"harkonnen/telemetry"
)

type Context struct {
	context.Context
	id              string
	variablePool    *VariablePool
	sampleCollector *telemetry.SampleCollector
	logger          *zerolog.Logger
	cancelFunc      context.CancelFunc
}

func NewContext(parent context.Context, parentLogger zerolog.Logger, shooterID string) Context {
	output := new(Context)
	output.sampleCollector = new(telemetry.SampleCollector)
	newLogger := parentLogger.With().Str("context", "Shooter").Str("ID", shooterID).Logger()
	output.logger = &newLogger
	output.id = shooterID
	output.variablePool = NewVariablePool(output.logger)
	output.Context, output.cancelFunc = context.WithCancel(parent)

	return *output
}

func (c *Context) Logger() *zerolog.Logger {
	return c.logger
}

func (c *Context) SampleCollector() *telemetry.SampleCollector {
	return c.sampleCollector
}

func (c *Context) VariablePool() *VariablePool {
	return c.variablePool
}

func (c *Context) ID() string {
	return c.id
}

func (c *Context) Cancel() {
	c.cancelFunc()
}

func (c *Context) NextLoop() <-chan struct{} {
	return make(chan struct{})
}

func (c *Context) OnUnrecoverableError(err error) {
	c.logger.Warn().Err(err).Msg("Caught an unrecoverable error, stopping current script execution")
	panic(err)
}
