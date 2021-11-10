package runtime_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/errors"
	"harkonnen/runtime"
	"testing"
)

func TestNewVariablePool(t *testing.T) {
	vp := runtime.NewVariablePool(&runtime.ErrorCollector{})
	assert.Implements(t, (*runtime.VariablesHandler)(nil), vp)
}

func TestVariablePool_SetGet(t *testing.T) {
	vp := runtime.NewVariablePool(&runtime.ErrorCollector{})
	vp.Set("test", "testValue")
	actualValue := vp.Get("test")

	assert.Equal(t, "testValue", actualValue)
}

func TestVariablePool_GetNonExisting(t *testing.T) {
	handler := &runtime.ErrorCollector{}
	vp := runtime.NewVariablePool(handler)
	actualValue := vp.Get("nonExisting")

	assert.Nil(t, actualValue, "Expected nil to be returned")

	collectedErrors := handler.GetCollected()

	assert.IsType(t, &errors.VariableNotFound{}, collectedErrors[0])
}

func TestVariablePool_GetString(t *testing.T) {
	vp := runtime.NewVariablePool(&runtime.ErrorCollector{})
	vp.Set("test", "anotherValue")
	actualValue := vp.Get("test")

	assert.Equal(t, "anotherValue", actualValue)
	assert.IsType(t, "", actualValue, "It is not a string")
}

func TestVariablePool_GetStringBadType(t *testing.T) {
	handler := &runtime.ErrorCollector{}
	vp := runtime.NewVariablePool(handler)
	vp.Set("test", 1)
	actualValue := vp.GetString("test")
	assert.Equal(t, "1", actualValue)

	collectedError := handler.GetCollected()
	assert.Nil(t, collectedError, "Expected no error when converting any type to string")
}

func TestVariablePool_GetInt(t *testing.T) {
	vp := runtime.NewVariablePool(&runtime.ErrorCollector{})
	vp.Set("test", 1)
	actualValue := vp.GetInt("test")

	assert.Equal(t, 1, actualValue)
}

func TestVariablePool_GetIntBadType(t *testing.T) {
	handler := &runtime.ErrorCollector{}
	vp := runtime.NewVariablePool(handler)
	vp.Set("test", "a")
	actualValue := vp.GetInt("test")

	assert.Equal(t, 0, actualValue, "Expected zero-valued int")

	collectedErrors := handler.GetCollected()
	assert.IsType(t, &errors.VariableCast{}, collectedErrors[0])
}

func TestVariablePool_GetBool(t *testing.T) {
	vp := runtime.NewVariablePool(&runtime.ErrorCollector{})
	vp.Set("test", true)
	actualValue := vp.GetBool("test")

	assert.Equal(t, true, actualValue)
}

func TestVariablePool_GetBoolBadType(t *testing.T) {
	handler := &runtime.ErrorCollector{}
	vp := runtime.NewVariablePool(handler)
	vp.Set("test", "booleanValue")
	actualValue := vp.GetBool("test")

	assert.Equal(t, false, actualValue, "Expected zero-valued bool (false)")

	collectedErrors := handler.GetCollected()
	assert.IsType(t, &errors.VariableCast{}, collectedErrors[0])
}

func TestVariablePool_Delete(t *testing.T) {
	handler := &runtime.ErrorCollector{}
	vp := runtime.NewVariablePool(handler)
	vp.Set("test", "finalValue")
	vp.Delete("test")
	actualResult := vp.Get("test")

	assert.Nil(t, actualResult)

	collectedErrors := handler.GetCollected()

	assert.IsType(t, &errors.VariableNotFound{}, collectedErrors[0])
}
