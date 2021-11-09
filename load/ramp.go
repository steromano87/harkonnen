package load

import (
	"math"
	"time"
)

type Ramp struct {
	loadUnits    int
	initialDelay time.Duration
	rampUpTime   time.Duration
	sustainTime  time.Duration
	rampDownTime time.Duration
}

func NewRamp(loadUnits int, initialDelay string, rampUpTime string, sustainTime string, rampDownTime string) (Ramp, error) {
	output := new(Ramp)
	var err error

	output.loadUnits = loadUnits

	output.initialDelay, err = time.ParseDuration(initialDelay)
	if err != nil {
		return Ramp{}, err
	}

	output.rampUpTime, err = time.ParseDuration(rampUpTime)
	if err != nil {
		return Ramp{}, err
	}

	output.sustainTime, err = time.ParseDuration(sustainTime)
	if err != nil {
		return Ramp{}, err
	}

	output.rampDownTime, err = time.ParseDuration(rampDownTime)
	if err != nil {
		return Ramp{}, err
	}

	return *output, nil
}

func (r Ramp) At(elapsed time.Duration) int {
	if elapsed < r.initialDelay {
		return 0
	}

	partialElapsed := elapsed - r.initialDelay

	if partialElapsed < r.rampUpTime {
		return int(
			math.Round(
				float64(r.loadUnits) * (float64(partialElapsed.Nanoseconds()) / float64(r.rampUpTime.Nanoseconds()))))
	}

	partialElapsed = partialElapsed - r.rampUpTime

	if partialElapsed < r.sustainTime {
		return r.loadUnits
	}

	partialElapsed = partialElapsed - r.sustainTime

	if partialElapsed < r.rampDownTime {
		return int(
			math.Round(
				float64(r.loadUnits) * float64(r.rampDownTime.Nanoseconds()-partialElapsed.Nanoseconds()) / float64(r.rampDownTime.Nanoseconds())))
	}

	return 0
}

func (r Ramp) TotalDuration() time.Duration {
	return r.initialDelay + r.rampUpTime + r.sustainTime + r.rampDownTime
}
