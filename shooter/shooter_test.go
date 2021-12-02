package shooter_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"harkonnen/shooter"
	"os"
	"sync"
	"testing"
)

type SampleClient struct {
	Context shooter.Context
}

func (c SampleClient) GenerateError() {
	c.Context.OnUnrecoverableError(errors.New("sample error"))
}

type ShooterTestSuite struct {
	suite.Suite
	logger    zerolog.Logger
	shooterID string

	setUpScript         shooter.Script
	mainScriptOne       shooter.Script
	mainScriptTwo       shooter.Script
	tearDownScript      shooter.Script
	explicitErrorScript shooter.Script
	implicitErrorScript shooter.Script
}

func (suite *ShooterTestSuite) SetupTest() {
	suite.shooterID = "1"
	suite.logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	suite.setUpScript = func(ctx shooter.Context) error {
		fmt.Println("Set Up script")
		return nil
	}

	suite.mainScriptOne = func(ctx shooter.Context) error {
		fmt.Println("Main script 1")
		return nil
	}

	suite.mainScriptTwo = func(ctx shooter.Context) error {
		fmt.Println("Main script 2")
		return nil
	}

	suite.tearDownScript = func(ctx shooter.Context) error {
		fmt.Println("Tear Down Script")
		return nil
	}

	suite.explicitErrorScript = func(ctx shooter.Context) error {
		return fmt.Errorf("triggered error")
	}

	suite.implicitErrorScript = func(ctx shooter.Context) error {
		client := SampleClient{Context: ctx}
		fmt.Println("Before error generation")
		client.GenerateError()
		fmt.Println("After error generation")
		return nil
	}
}

func (suite *ShooterTestSuite) createContext() shooter.Context {
	return shooter.NewContext(context.Background(), suite.logger, suite.shooterID)
}

func (suite *ShooterTestSuite) TestStart() {
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		Context:        suite.createContext(),
		SetUpScript:    suite.setUpScript,
		MainScripts:    []shooter.Script{suite.mainScriptOne, suite.mainScriptTwo},
		TearDownScript: suite.tearDownScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(suite.T(), shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(suite.T(), 3, testShooter.TotalIterations())
	assert.Equal(suite.T(), testShooter.TotalIterations(), testShooter.SuccessfulIterations())
	assert.Equal(suite.T(), shooter.Completed, testShooter.Status())
}

func (suite *ShooterTestSuite) TestExplicitErrorInMainLoop() {
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		Context:        suite.createContext(),
		SetUpScript:    suite.setUpScript,
		MainScripts:    []shooter.Script{suite.mainScriptOne, suite.explicitErrorScript},
		TearDownScript: suite.tearDownScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(suite.T(), shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(suite.T(), 3, testShooter.TotalIterations())
	assert.Equal(suite.T(), 0, testShooter.SuccessfulIterations())
	assert.Equal(suite.T(), shooter.Completed, testShooter.Status())
}

func (suite *ShooterTestSuite) TestImplicitErrorInMainLoop() {
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		Context:        suite.createContext(),
		SetUpScript:    nil,
		MainScripts:    []shooter.Script{suite.implicitErrorScript, suite.mainScriptTwo},
		TearDownScript: nil,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(suite.T(), shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(suite.T(), 3, testShooter.TotalIterations())
	assert.Equal(suite.T(), 0, testShooter.SuccessfulIterations())
	assert.Equal(suite.T(), shooter.Completed, testShooter.Status())
}

func (suite *ShooterTestSuite) TestImplicitErrorInSetupScript() {
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		Context:        suite.createContext(),
		SetUpScript:    suite.implicitErrorScript,
		MainScripts:    []shooter.Script{suite.mainScriptOne, suite.mainScriptTwo},
		TearDownScript: nil,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(suite.T(), shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(suite.T(), 0, testShooter.TotalIterations())
	assert.Equal(suite.T(), 0, testShooter.SuccessfulIterations())
	assert.Equal(suite.T(), shooter.Error, testShooter.Status())
}

func (suite *ShooterTestSuite) TestImplicitErrorInTearDownScript() {
	wg := sync.WaitGroup{}

	testShooter := shooter.Shooter{
		Context:        suite.createContext(),
		SetUpScript:    nil,
		MainScripts:    []shooter.Script{suite.mainScriptOne, suite.mainScriptTwo},
		TearDownScript: suite.implicitErrorScript,
		MaxIterations:  3,
		WaitGroup:      &wg,
	}

	wg.Add(1)
	testShooter.Start()
	assert.Equal(suite.T(), shooter.Running, testShooter.Status())

	wg.Wait()
	assert.Equal(suite.T(), 3, testShooter.TotalIterations())
	assert.Equal(suite.T(), 3, testShooter.SuccessfulIterations())
	assert.Equal(suite.T(), shooter.Error, testShooter.Status())
}

func TestShooterTestSuite(t *testing.T) {
	suite.Run(t, new(ShooterTestSuite))
}
