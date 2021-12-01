package shooter

import (
	"github.com/rs/zerolog"
	"harkonnen/telemetry"
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

	SampleCollector *telemetry.SampleCollector
	VariablePool    *VariablePool

	status               Status
	totalIterations      int
	successfulIterations int
	scheduledForShutdown bool
}

func (s *Shooter) Start() {
	go s.run()
	s.status = Running
	s.Logger().Info().Msg("Shooter started")
}

func (s *Shooter) ScheduleShutDown() {
	s.Logger().Info().Msg("Shooter marked for shutdown")
	s.scheduledForShutdown = true
	s.status = ShuttingDown
}

func (s *Shooter) TotalIterations() int {
	return s.totalIterations
}

func (s *Shooter) SuccessfulIterations() int {
	return s.successfulIterations
}

func (s *Shooter) Logger() *zerolog.Logger {
	shooterLogger := s.Context.Logger().With().Str("ID", s.ID).Logger()
	return &shooterLogger
}

func (s *Shooter) run() {
	defer s.WaitGroup.Done()

	// Setup script execution
	s.executeSetupScript()
	if s.status == Error {
		return
	}

	// Main loop
	s.executeMainScripts()

	// Teardown script execution
	s.executeTearDownScript()
	if s.status == Error {
		return
	}

	if s.scheduledForShutdown {
		s.status = Stopped
	} else {
		s.status = Completed
	}
}

func (s *Shooter) executeSetupScript() {
	if s.SetUpScript != nil {
		defer s.handleSetUpTearDownPanic()

		s.Logger().Info().Msg("Started setup script execution")
		err := s.SetUpScript(s.Context)
		s.Logger().Info().Msg("Setup script execution completed")

		if err != nil {
			s.Context.OnUnrecoverableError(err)
		}
	}
}

func (s *Shooter) executeMainScripts() {
	if len(s.MainScripts) > 0 {
		for !s.scheduledForShutdown && (s.totalIterations < s.MaxIterations || s.MaxIterations == 0) {
			s.executeMainScriptsLoop()
		}
	}
}

func (s *Shooter) executeMainScriptsLoop() {
	defer s.handleMainLoopPanic()

	var err error

	for _, mainScript := range s.MainScripts {
		// Check for termination at every loop
		select {
		case <-s.Context.Done():
			s.status = ShuttingDown
			return

		case <-s.Context.NextLoop():
			break

		default:
		}

		err := mainScript(s.Context)
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

func (s *Shooter) executeTearDownScript() {
	if s.TearDownScript != nil {
		defer s.handleSetUpTearDownPanic()

		s.Logger().Info().Msg("Started teardown script execution")
		err := s.TearDownScript(s.Context)
		s.Logger().Info().Msg("Teardown script execution completed")

		if err != nil {
			s.Context.OnUnrecoverableError(err)
		}
	}
}

func (s *Shooter) handleSetUpTearDownPanic() {
	if err := recover(); err != nil {
		s.Logger().Error().Stack().Err(err.(error)).Msg("Encountered error during setup/teardown sequence")
		s.status = Error
		return
	}
}

func (s *Shooter) handleMainLoopPanic() {
	if err := recover(); err != nil {
		s.Logger().Error().Stack().Err(err.(error)).Msg("Encountered error during main loop")
		s.totalIterations++
		return
	}
}

func (s *Shooter) Status() Status {
	return s.status
}
