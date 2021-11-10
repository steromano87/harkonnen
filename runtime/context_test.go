package runtime_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"harkonnen/runtime"
	"testing"
)

func TestNewContext(t *testing.T) {
	errorHandler := runtime.ErrorCollector{}
	testContext := runtime.NewContext(context.Background(), runtime.NewVariablePool(&errorHandler), &errorHandler, runtime.NewSampleCollector())

	assert.IsType(t, &runtime.Context{}, testContext)
	assert.Implements(t, (*runtime.VariablesHandler)(nil), testContext)
}
