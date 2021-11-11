package load_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"harkonnen/load"
	"harkonnen/runtime"
	"sync"
	"testing"
)

var setUpScript = func(ctx runtime.Context) error {
	fmt.Println("Set Up script")
	return nil
}

var mainScriptOne = func(ctx runtime.Context) error {
	fmt.Println("Main script 1")
	return nil
}

var mainScriptTwo = func(ctx runtime.Context) error {
	fmt.Println("Main script 2")
	return nil
}

var tearDownScript = func(ctx runtime.Context) error {
	fmt.Println("Tear Down Script")
	return nil
}

var scriptWithError = func(ctx runtime.Context) error {
	return fmt.Errorf("triggered error")
}

func TestUnit_Start(t *testing.T) {
	errorsCollector := runtime.ErrorCollector{}
	testContext := runtime.NewContext(context.Background(), runtime.NewVariablePool(&errorsCollector), &errorsCollector, runtime.NewSampleCollector())
	wg := sync.WaitGroup{}

	testUnit := load.Unit{
		ID:             "Sample",
		Context:        testContext,
		SetUpScript:    setUpScript,
		MainScripts:    []load.Script{mainScriptOne, mainScriptTwo},
		TearDownScript: tearDownScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testUnit.Start()
	assert.Equal(t, load.Running, testUnit.Status())

	wg.Wait()
	assert.Equal(t, 3, testUnit.TotalIterations())
	assert.Equal(t, testUnit.TotalIterations(), testUnit.SuccessfulIterations())
	assert.Equal(t, load.Completed, testUnit.Status())
}

func TestUnit_StartWithError(t *testing.T) {
	errorsCollector := runtime.ErrorCollector{}
	testContext := runtime.NewContext(context.Background(), runtime.NewVariablePool(&errorsCollector), &errorsCollector, runtime.NewSampleCollector())
	wg := sync.WaitGroup{}

	testUnit := load.Unit{
		ID:             "Sample",
		Context:        testContext,
		SetUpScript:    setUpScript,
		MainScripts:    []load.Script{mainScriptOne, scriptWithError},
		TearDownScript: tearDownScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testUnit.Start()
	assert.Equal(t, load.Running, testUnit.Status())

	wg.Wait()
	assert.Equal(t, 3, testUnit.TotalIterations())
	assert.Equal(t, 0, testUnit.SuccessfulIterations())
	assert.True(t, testContext.HasErrors())
	assert.Equal(t, load.Completed, testUnit.Status())
}
