package shooter_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"harkonnen/shooter"
	"sync"
	"testing"
)

type SampleClient struct {
	Context shooter.Context
}

func (c SampleClient) GenerateError() {
	c.Context.OnUnrecoverableError(errors.New("sample error"))
}

var setUpScript = func(ctx shooter.Context) error {
	fmt.Println("Set Up script")
	return nil
}

var mainScriptOne = func(ctx shooter.Context) error {
	fmt.Println("Main script 1")
	return nil
}

var mainScriptTwo = func(ctx shooter.Context) error {
	fmt.Println("Main script 2")
	return nil
}

var tearDownScript = func(ctx shooter.Context) error {
	fmt.Println("Tear Down Script")
	return nil
}

var scriptWithError = func(ctx shooter.Context) error {
	return fmt.Errorf("triggered error")
}

var scriptWithImplicitError = func(ctx shooter.Context) error {
	client := SampleClient{Context: ctx}
	fmt.Println("Before error generation")
	client.GenerateError()
	fmt.Println("After error generation")
	return nil
}

func TestShooter_Start(t *testing.T) {
	testContext := shooter.NewContext(context.Background())
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		ID:             "BaseSample",
		Context:        testContext,
		SetUpScript:    setUpScript,
		MainScripts:    []shooter.Script{mainScriptOne, mainScriptTwo},
		TearDownScript: tearDownScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(t, shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(t, 3, testShooter.TotalIterations())
	assert.Equal(t, testShooter.TotalIterations(), testShooter.SuccessfulIterations())
	assert.Equal(t, shooter.Completed, testShooter.Status())
}

func TestShooter_StartWithError(t *testing.T) {
	testContext := shooter.NewContext(context.Background())
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		ID:             "BaseSample",
		Context:        testContext,
		SetUpScript:    setUpScript,
		MainScripts:    []shooter.Script{mainScriptOne, scriptWithError},
		TearDownScript: tearDownScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(t, shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(t, 3, testShooter.TotalIterations())
	assert.Equal(t, 0, testShooter.SuccessfulIterations())
	assert.Equal(t, shooter.Completed, testShooter.Status())
}

func TestShooter_ScriptWithImplicitError(t *testing.T) {
	testContext := shooter.NewContext(context.Background())
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		ID:             "BaseShooter",
		Context:        testContext,
		SetUpScript:    nil,
		MainScripts:    []shooter.Script{scriptWithImplicitError, mainScriptTwo},
		TearDownScript: nil,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(t, shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(t, 3, testShooter.TotalIterations())
	assert.Equal(t, 0, testShooter.SuccessfulIterations())
	assert.Equal(t, shooter.Completed, testShooter.Status())
}
