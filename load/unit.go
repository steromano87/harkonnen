package load

import (
	"harkonnen/runtime"
	"sync"
)

type Unit struct {
	ID             string
	Context        runtime.Context
	SetUpScript    Script
	MainScripts    []Script
	TearDownScript Script
	MaxIterations  int
	WaitGroup      *sync.WaitGroup

	status               UnitStatus
	totalIterations      int
	successfulIterations int
	scheduledForShutdown bool
}

func (u *Unit) Start() {
	u.status = Running
	go u.run()
}

func (u *Unit) ScheduleShutDown() {
	u.scheduledForShutdown = true
	u.status = ShuttingDown
}

func (u *Unit) TotalIterations() int {
	return u.totalIterations
}

func (u *Unit) SuccessfulIterations() int {
	return u.successfulIterations
}

func (u *Unit) run() {
	defer u.WaitGroup.Done()

	// Setup script execution
	err := u.executeSetupScript()
	if err != nil {
		u.status = Error
		return
	}

	// Main loop
	u.executeMainScripts()

	// Teardown script execution
	err = u.executeTearDownScript()
	if err != nil {
		u.status = Error
		return
	}
}

func (u *Unit) executeSetupScript() error {
	if u.SetUpScript != nil {
		return u.SetUpScript(u.Context)
	}

	return nil
}

func (u *Unit) executeMainScripts() {
	if len(u.MainScripts) > 0 {
		for !u.scheduledForShutdown && (u.totalIterations < u.MaxIterations || u.MaxIterations == 0) {
			var err error

			for _, script := range u.MainScripts {
				// Check for termination at every loop
				select {
				case <-u.Context.Done():
					u.status = ShuttingDown
					return

				case <-u.Context.NextLoop():
					break

				default:
				}

				err = script(u.Context)
				// Exit from inner loop in case of error while executing one of the scripts
				if err != nil {
					u.Context.Capture(err)
					u.Context.NextLoop()
				}
			}

			u.totalIterations++

			if err == nil {
				u.successfulIterations++
			}
		}
	}
}

func (u *Unit) executeTearDownScript() error {
	if u.TearDownScript != nil {
		return u.TearDownScript(u.Context)
	}

	return nil
}

func (u *Unit) Status() UnitStatus {
	return u.status
}
