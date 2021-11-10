package load_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"harkonnen/load"
	"harkonnen/runtime"
	"testing"
)

func TestNewContext(t *testing.T) {
	errorHandler := runtime.ErrorHandler{}
	testContext := load.NewContext(context.Background(), runtime.NewVariablePool(&errorHandler), &errorHandler, runtime.NewSampleCollector())

	assert.IsType(t, &load.Context{}, testContext)
	assert.Implements(t, (*runtime.VariablesHandler)(nil), testContext)
}
