package shooter_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/log"
	"harkonnen/shooter"
	"testing"
)

func TestVariablePool_SetGet(t *testing.T) {
	vp := shooter.NewVariablePool(&log.Collector{})
	vp.Set("test", "testValue")
	actualValue, err := vp.Get("test")

	assert.Equal(t, "testValue", actualValue)
	assert.NoError(t, err)
}

func TestVariablePool_GetNonExisting(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	actualValue, err := vp.Get("nonExisting")

	assert.Nil(t, actualValue, "Expected nil to be returned")
	assert.IsType(t, shooter.ErrVariableNotFound{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}

func TestVariablePool_GetString(t *testing.T) {
	vp := shooter.NewVariablePool(&log.Collector{})
	vp.Set("test", "anotherValue")
	actualValue, err := vp.Get("test")

	assert.Equal(t, "anotherValue", actualValue)
	assert.IsType(t, "", actualValue, "It is not a string")
	assert.NoError(t, err)
}

func TestVariablePool_GetStringNonExisting(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", "anotherValue")
	actualValue, err := vp.GetInt("nonExisting")

	assert.Empty(t, actualValue, "Expected empty to be returned")
	assert.IsType(t, shooter.ErrVariableNotFound{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}

func TestVariablePool_GetStringBadType(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", 1)
	actualValue, err := vp.GetString("test")
	assert.Equal(t, "1", actualValue)
	assert.NoError(t, err)
	assert.Empty(t, logCollector.Flush(log.ErrorLevel))
}

func TestVariablePool_GetInt(t *testing.T) {
	vp := shooter.NewVariablePool(&log.Collector{})
	vp.Set("test", 1)
	actualValue, err := vp.GetInt("test")

	assert.Equal(t, 1, actualValue)
	assert.NoError(t, err)
}

func TestVariablePool_GetIntNonExisting(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", 7)
	actualValue, err := vp.GetInt("nonExisting")

	assert.Zero(t, actualValue, "Expected zero to be returned")
	assert.IsType(t, shooter.ErrVariableNotFound{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}

func TestVariablePool_GetIntBadType(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", "a")
	actualValue, err := vp.GetInt("test")

	assert.Equal(t, 0, actualValue, "Expected zero-valued int")
	assert.IsType(t, shooter.ErrBadVariableCast{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}

func TestVariablePool_GetBool(t *testing.T) {
	vp := shooter.NewVariablePool(&log.Collector{})
	vp.Set("test", true)
	actualValue, err := vp.GetBool("test")

	assert.Equal(t, true, actualValue)
	assert.NoError(t, err)
}

func TestVariablePool_GetBoolNonExisting(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", true)
	actualValue, err := vp.GetBool("nonExisting")

	assert.False(t, actualValue, "Expected false to be returned")
	assert.IsType(t, shooter.ErrVariableNotFound{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}

func TestVariablePool_GetBoolBadType(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", "booleanValue")
	actualValue, err := vp.GetBool("test")

	assert.Equal(t, false, actualValue, "Expected zero-valued bool (false)")
	assert.IsType(t, shooter.ErrBadVariableCast{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}

func TestVariablePool_Delete(t *testing.T) {
	logCollector := log.Collector{}
	vp := shooter.NewVariablePool(&logCollector)
	vp.Set("test", "finalValue")
	vp.Delete("test")
	actualResult, err := vp.Get("test")

	assert.Nil(t, actualResult)
	assert.IsType(t, shooter.ErrVariableNotFound{}, err)
	assert.Equal(t, 1, len(logCollector.Flush(log.ErrorLevel)))
}
