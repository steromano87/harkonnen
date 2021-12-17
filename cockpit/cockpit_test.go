package cockpit_test

import (
	"github.com/steromano87/harkonnen/cockpit"
	"github.com/steromano87/harkonnen/injector"
	"github.com/steromano87/harkonnen/load"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
	"time"
)

type CockpitTestSuite struct {
	suite.Suite
	loadRampOdd     load.LinearRamp
	loadRampEven    load.LinearRamp
	injectorEven    injector.Reference
	injectorOdd     injector.Reference
	standardElapsed time.Duration
}

func (suite *CockpitTestSuite) SetupTest() {
	suite.loadRampOdd, _ = load.ParseLinearRamp(5, "0s", "1s", "10s", "1s")
	suite.loadRampEven, _ = load.ParseLinearRamp(6, "0s", "1s", "10s", "1s")

	suite.injectorOdd = injector.Reference{
		Address: "8.8.8.8",
		Port:    3200,
		Weight:  1,
		Labels:  nil,
		Type:    injector.LocalInjector,
	}

	suite.injectorEven = injector.Reference{
		Address: "8.8.8.8",
		Port:    3200,
		Weight:  2,
		Labels:  nil,
		Type:    injector.LocalInjector,
	}

	suite.standardElapsed, _ = time.ParseDuration("6s")
}

func (suite *CockpitTestSuite) TestAtForEachInjector_UniformWeightsWithoutRemainder() {
	cc := cockpit.Cockpit{
		Context:      nil,
		InjectorType: injector.LocalInjector,
		Injectors:    map[string]injector.Reference{"even1": suite.injectorEven, "even2": suite.injectorEven},
		LoadProfiles: []load.Profile{suite.loadRampEven},
	}

	assert.Equal(suite.T(), 6, cc.At(suite.standardElapsed))
	shooterQuotas := cc.AtForEachInjector(suite.standardElapsed)

	assert.Equal(suite.T(), 3, shooterQuotas["even1"])
	assert.Equal(suite.T(), 3, shooterQuotas["even2"])
}

func (suite CockpitTestSuite) TestAtForEachInjector_UniformWeightsWithRemainder() {
	cc := cockpit.Cockpit{
		Context:      nil,
		InjectorType: injector.LocalInjector,
		Injectors:    map[string]injector.Reference{"even1": suite.injectorEven, "even2": suite.injectorEven},
		LoadProfiles: []load.Profile{suite.loadRampOdd},
	}

	assert.Equal(suite.T(), 5, cc.At(suite.standardElapsed))
	shooterQuotas := cc.AtForEachInjector(suite.standardElapsed)

	// Since assignation is sometimes random, we can only calculate that:
	// - the sum of the quotas equals the total amount of shooters to distribute
	// - there is only a difference of 1 between the max and min quota
	totalCalculatedShooters := 0
	minQuota := 99
	maxQuota := 0
	for _, quota := range shooterQuotas {
		totalCalculatedShooters += quota
		if quota < minQuota {
			minQuota = quota
		}
		if quota > maxQuota {
			maxQuota = quota
		}
	}

	assert.Equal(suite.T(), 5, totalCalculatedShooters, "total of quotas differs from the amount of shooters")
	assert.Equal(suite.T(), 1, maxQuota-minQuota, "incorrect distribution of quotas")
}

func (suite CockpitTestSuite) TestAtForEachInjector_UniformWeightsWithRemainderMoreInjectors() {
	cc := cockpit.Cockpit{
		Context:      nil,
		InjectorType: injector.LocalInjector,
		Injectors:    map[string]injector.Reference{"even1": suite.injectorEven, "even2": suite.injectorEven, "even3": suite.injectorEven},
		LoadProfiles: []load.Profile{suite.loadRampOdd},
	}

	assert.Equal(suite.T(), 5, cc.At(suite.standardElapsed))
	shooterQuotas := cc.AtForEachInjector(suite.standardElapsed)

	// Since assignation is sometimes random, we can only calculate that:
	// - the sum of the quotas equals the total amount of shooters to distribute
	// - there is only a difference of 1 between the max and min quota
	totalCalculatedShooters := 0
	minQuota := 99
	maxQuota := 0
	for _, quota := range shooterQuotas {
		totalCalculatedShooters += quota
		if quota < minQuota {
			minQuota = quota
		}
		if quota > maxQuota {
			maxQuota = quota
		}
	}

	assert.Equal(suite.T(), 5, totalCalculatedShooters, "total of quotas differs from the amount of shooters")
	assert.Equal(suite.T(), 1, maxQuota-minQuota, "incorrect distribution of quotas")
}

func (suite CockpitTestSuite) TestAtForEachInjector_NonUniformWeightsWithoutRemainder() {
	cc := cockpit.Cockpit{
		Context:      nil,
		InjectorType: injector.LocalInjector,
		Injectors:    map[string]injector.Reference{"even1": suite.injectorEven, "even2": suite.injectorEven, "odd1": suite.injectorOdd},
		LoadProfiles: []load.Profile{suite.loadRampOdd},
	}

	assert.Equal(suite.T(), 5, cc.At(suite.standardElapsed))
	shooterQuotas := cc.AtForEachInjector(suite.standardElapsed)

	assert.Equal(suite.T(), 2, shooterQuotas["even1"])
	assert.Equal(suite.T(), 2, shooterQuotas["even2"])
	assert.Equal(suite.T(), 1, shooterQuotas["odd1"])
}

func (suite CockpitTestSuite) TestAtForEachInjector_NonUniformWeightsWithRemainder() {
	cc := cockpit.Cockpit{
		Context:      nil,
		InjectorType: injector.LocalInjector,
		Injectors:    map[string]injector.Reference{"even1": suite.injectorEven, "odd1": suite.injectorOdd},
		LoadProfiles: []load.Profile{suite.loadRampOdd, suite.loadRampEven},
	}

	assert.Equal(suite.T(), 11, cc.At(suite.standardElapsed))
	shooterQuotas := cc.AtForEachInjector(suite.standardElapsed)

	// Since assignation is sometimes random, we can only calculate that:
	// - the sum of the quotas equals the total amount of shooters to distribute
	// - the quota ratio is APPROXIMATELY 2:1 (it will be rounded to the nearest integer
	totalCalculatedShooters := 0
	minQuota := 99
	maxQuota := 0
	for _, quota := range shooterQuotas {
		totalCalculatedShooters += quota
		if quota < minQuota {
			minQuota = quota
		}
		if quota > maxQuota {
			maxQuota = quota
		}
	}

	assert.Equal(suite.T(), 11, totalCalculatedShooters, "total of quotas differs from the amount of shooters")
	assert.Equal(suite.T(), 2, int(math.Round(float64(maxQuota)/float64(minQuota))))
}

func TestCockpitTestSuite(t *testing.T) {
	suite.Run(t, new(CockpitTestSuite))
}
