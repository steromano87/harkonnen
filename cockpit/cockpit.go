package cockpit

import (
	"context"
	"github.com/steromano87/harkonnen/injector"
	"github.com/steromano87/harkonnen/load"
	"math"
	"time"
)

type Cockpit struct {
	Context context.Context

	InjectorType injector.Type
	Injectors    map[string]injector.Reference

	LoadProfiles []load.Profile
}

func (c Cockpit) At(elapsed time.Duration) int {
	totalInjectors := 0
	for _, loadProfile := range c.LoadProfiles {
		totalInjectors += loadProfile.At(elapsed)
	}

	return totalInjectors
}

func (c Cockpit) TotalDuration() time.Duration {
	maxDuration, _ := time.ParseDuration("0s")
	for _, loadProfile := range c.LoadProfiles {
		if loadProfile.TotalDuration() > maxDuration {
			maxDuration = loadProfile.TotalDuration()
		}
	}

	return maxDuration
}

func (c Cockpit) AtForEachInjector(elapsed time.Duration) map[string]int {
	totalShooters := c.At(elapsed)

	totalWeights := 0
	for _, spec := range c.Injectors {
		if spec.Weight == 0 {
			totalWeights += 1
		} else {
			totalWeights += spec.Weight
		}
	}

	output := make(map[string]int)

	for id, spec := range c.Injectors {
		shootersQuota := int(math.Round(float64(totalShooters*spec.Weight) / float64(totalWeights)))

		output[id] = shootersQuota
		totalWeights -= spec.Weight
		totalShooters -= shootersQuota
	}
	return output
}
