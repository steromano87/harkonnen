package load_test

import (
	"fmt"
	"github.com/steromano87/harkonnen/load"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
	"time"
)

var shooters = 8
var initialDelay = "5s"
var rampUpTime = "4s"
var sustainTime = "10s"
var rampDownTime = "2s"
var testRamp, _ = load.ParseLinearRamp(shooters, initialDelay, rampUpTime, sustainTime, rampDownTime)

func TestRampCreation(t *testing.T) {
	assert.IsType(t, load.LinearRamp{}, testRamp)
}

func TestRamp_At_BeforeStart(t *testing.T) {
	elapsed, _ := time.ParseDuration("-2s")
	assert.Equal(t, 0, testRamp.At(elapsed))
}

func TestRamp_At_DuringInitialDelay(t *testing.T) {
	elapsed, _ := time.ParseDuration("3s")
	assert.Equal(t, 0, testRamp.At(elapsed))
}

func TestRamp_At_DuringRampUp(t *testing.T) {
	elapsed, _ := time.ParseDuration("6s")
	assert.Equal(t, 2, testRamp.At(elapsed))
}

func TestRamp_At_DuringSustain(t *testing.T) {
	elapsed, _ := time.ParseDuration("12s")
	assert.Equal(t, 8, testRamp.At(elapsed))
}

func TestRamp_At_DuringRampDown(t *testing.T) {
	elapsed, _ := time.ParseDuration("20s")
	assert.Equal(t, 4, testRamp.At(elapsed))
}

func TestRamp_At_AfterCompletion(t *testing.T) {
	elapsed, _ := time.ParseDuration("30s")
	assert.Equal(t, 0, testRamp.At(elapsed))
}

func TestRamp_TotalDuration(t *testing.T) {
	elapsed, _ := time.ParseDuration("21s")
	assert.Equal(t, elapsed, testRamp.TotalDuration())
}

func TestRamp_Marshalling(t *testing.T) {
	output, err := yaml.Marshal(testRamp)

	if assert.NoError(t, err) {
		assert.Contains(t, string(output), fmt.Sprintf("shooters: %d", shooters))
		assert.Contains(t, string(output), fmt.Sprintf("initial_delay: %s", initialDelay))
		assert.Contains(t, string(output), fmt.Sprintf("ramp_up_time: %s", rampUpTime))
		assert.Contains(t, string(output), fmt.Sprintf("sustain_time: %s", sustainTime))
		assert.Contains(t, string(output), fmt.Sprintf("ramp_down_time: %s", rampDownTime))
	}
}

func TestRamp_UnmarshalYAML(t *testing.T) {
	output, _ := yaml.Marshal(testRamp)

	var outputRamp load.LinearRamp
	err := yaml.Unmarshal(output, &outputRamp)

	if assert.NoError(t, err) {
		assert.Equal(t, shooters, outputRamp.Shooters)

		parsedInitialDelay, _ := time.ParseDuration(initialDelay)
		parsedRampUpTime, _ := time.ParseDuration(rampUpTime)
		parsedSustainTime, _ := time.ParseDuration(sustainTime)
		parsedRampDownTime, _ := time.ParseDuration(rampDownTime)

		assert.Equal(t, parsedInitialDelay, outputRamp.InitialDelay)
		assert.Equal(t, parsedRampUpTime, outputRamp.RampUpTime)
		assert.Equal(t, parsedSustainTime, outputRamp.SustainTime)
		assert.Equal(t, parsedRampDownTime, outputRamp.RampDownTime)
	}

}
