package shooter

import (
	"fmt"
	"github.com/go-errors/errors"
	"sync"
)

type Shooter struct {
	ID             string
	Context        Context
	SetUpScript    Script
	MainScripts    []Script
	TearDownScript Script
	MaxIterations  int
	WaitGroup      *sync.WaitGroup

	status               Status
	totalIterations      int
	successfulIterations int
	scheduledForShutdown bool
}

func (s *Shooter) Start() {
	s.status = Running
	s.Context.LogCollector().Info(fmt.Sprintf("Shooter %s started", s.ID))
	go s.run()
}

func (s *Shooter) ScheduleShutDown() {
	s.Context.LogCollector().Info(fmt.Sprintf("Shooter %s marked for shotdown", s.ID))
	s.scheduledForShutdown = true
	s.status = ShuttingDown
}

func (s *Shooter) TotalIterations() int {
	return s.totalIterations
}

func (s *Shooter) SuccessfulIterations() int {
	return s.successfulIterations
}

func (s *Shooter) run() {
	defer s.WaitGroup.Done()

	// Setup script execution
	err := s.executeSetupScript()
	if err != nil {
		s.status = Error
		return
	}

	// Main loop
	s.executeMainScripts()

	// Teardown script execution
	err = s.executeTearDownScript()
	if err != nil {
		s.status = Error
		return
	}

	if s.scheduledForShutdown {
		s.status = Stopped
	} else {
		s.status = Completed
	}
}

func (s *Shooter) executeSetupScript() error {
	if s.SetUpScript != nil {
		defer s.handleLoopPanic()

		s.Context.LogCollector().Info("Started setup script execution")
		err := s.SetUpScript(s.Context)
		s.Context.LogCollector().Info("Setup script execution completed")

		if err != nil {
			s.Context.LogCollector().Error(err.Error())
		}

		return err
	}

	return nil
}

func (s *Shooter) executeMainScripts() {
	if len(s.MainScripts) > 0 {
		for !s.scheduledForShutdown && (s.totalIterations < s.MaxIterations || s.MaxIterations == 0) {
			s.executeMainScriptsLoop()
		}
	}
}

func (s *Shooter) executeMainScriptsLoop() {
	defer s.handleLoopPanic()

	var err error

	for _, script := range s.MainScripts {
		// Check for termination at every loop
		select {
		case <-s.Context.Done():
			s.status = ShuttingDown
			return

		case <-s.Context.NextLoop():
			break

		default:
		}

		err := script(s.Context)
		// Exit from inner loop in case of error while executing one of the scripts
		if err != nil {
			s.Context.OnUnrecoverableError(err)
		}
	}

	s.totalIterations++

	if err == nil {
		s.successfulIterations++
	}
}

func (s *Shooter) executeTearDownScript() error {
	if s.TearDownScript != nil {
		defer s.handleLoopPanic()

		s.Context.LogCollector().Info("Started teardown script execution")
		err := s.TearDownScript(s.Context)
		s.Context.LogCollector().Info("Teardown script execution completed")

		if err != nil {
			s.Context.LogCollector().Error(err.Error())
		}
	}

	return nil
}

func (s *Shooter) handleLoopPanic() {
	if err := recover(); err != nil {
		s.Context.LogCollector().Error(errors.Wrap(err, 2).ErrorStack())
		s.totalIterations++
		return
	}
}

func (s *Shooter) Status() Status {
	return s.status
}
