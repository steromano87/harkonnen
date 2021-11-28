package shooter

import (
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
	go s.run()
}

func (s *Shooter) ScheduleShutDown() {
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
		return s.SetUpScript(s.Context)
	}

	return nil
}

func (s *Shooter) executeMainScripts() {
	if len(s.MainScripts) > 0 {
		for !s.scheduledForShutdown && (s.totalIterations < s.MaxIterations || s.MaxIterations == 0) {
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

				err = script(s.Context)
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
	}
}

func (s *Shooter) executeTearDownScript() error {
	if s.TearDownScript != nil {
		return s.TearDownScript(s.Context)
	}

	return nil
}

func (s *Shooter) Status() Status {
	return s.status
}
