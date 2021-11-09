package load_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/load"
	"testing"
	"time"
)

var testRamp, _ = load.NewRamp(8, "5s", "4s", "10s", "2s")

func TestRampCreation(t *testing.T) {
	assert.IsType(t, load.Ramp{}, testRamp)
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
