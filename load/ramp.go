package load

import (
	"gopkg.in/yaml.v3"
	"math"
	"time"
)

type Ramp struct {
	LoadUnits    int           `yaml:"load_units"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	RampUpTime   time.Duration `yaml:"ramp_up_time"`
	SustainTime  time.Duration `yaml:"sustain_time"`
	RampDownTime time.Duration `yaml:"ramp_down_time"`
}

func NewRamp(loadUnits int, initialDelay string, rampUpTime string, sustainTime string, rampDownTime string) (Ramp, error) {
	output := new(Ramp)
	var err error

	output.LoadUnits = loadUnits

	output.InitialDelay, err = time.ParseDuration(initialDelay)
	if err != nil {
		return Ramp{}, err
	}

	output.RampUpTime, err = time.ParseDuration(rampUpTime)
	if err != nil {
		return Ramp{}, err
	}

	output.SustainTime, err = time.ParseDuration(sustainTime)
	if err != nil {
		return Ramp{}, err
	}

	output.RampDownTime, err = time.ParseDuration(rampDownTime)
	if err != nil {
		return Ramp{}, err
	}

	return *output, nil
}

func (r *Ramp) At(elapsed time.Duration) int {
	if elapsed < r.InitialDelay {
		return 0
	}

	partialElapsed := elapsed - r.InitialDelay

	if partialElapsed < r.RampUpTime {
		return int(
			math.Round(
				float64(r.LoadUnits) * (float64(partialElapsed.Nanoseconds()) / float64(r.RampUpTime.Nanoseconds()))))
	}

	partialElapsed = partialElapsed - r.RampUpTime

	if partialElapsed < r.SustainTime {
		return r.LoadUnits
	}

	partialElapsed = partialElapsed - r.SustainTime

	if partialElapsed < r.RampDownTime {
		return int(
			math.Round(
				float64(r.LoadUnits) * float64(r.RampDownTime.Nanoseconds()-partialElapsed.Nanoseconds()) / float64(r.RampDownTime.Nanoseconds())))
	}

	return 0
}

func (r *Ramp) TotalDuration() time.Duration {
	return r.InitialDelay + r.RampUpTime + r.SustainTime + r.RampDownTime
}

func (r *Ramp) UnmarshalYAML(value *yaml.Node) error {
	var temp struct {
		LoadUnits    int    `yaml:"load_units"`
		InitialDelay string `yaml:"initial_delay"`
		RampUpTime   string `yaml:"ramp_up_time"`
		SustainTime  string `yaml:"sustain_time"`
		RampDownTime string `yaml:"ramp_down_time"`
	}

	if err := value.Decode(&temp); err != nil {
		return err
	}

	var err error
	r.LoadUnits = temp.LoadUnits
	r.InitialDelay, err = time.ParseDuration(temp.InitialDelay)
	if err != nil {
		return err
	}

	r.RampUpTime, err = time.ParseDuration(temp.RampUpTime)
	if err != nil {
		return err
	}

	r.SustainTime, err = time.ParseDuration(temp.SustainTime)
	if err != nil {
		return err
	}

	r.RampDownTime, err = time.ParseDuration(temp.RampDownTime)
	if err != nil {
		return err
	}

	return nil
}
