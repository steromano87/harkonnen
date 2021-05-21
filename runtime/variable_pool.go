package runtime

import (
	"fmt"
	"harkonnen/errors"
)

type VariablePool struct {
	variables    map[string]interface{}
	errorHandler *ErrorHandler
}

func NewVariablePool(errorHandler *ErrorHandler) *VariablePool {
	pool := new(VariablePool)
	pool.variables = make(map[string]interface{})
	pool.errorHandler = errorHandler

	return pool
}

func (pool *VariablePool) Set(name string, value interface{}) {
	pool.variables[name] = value
}

func (pool *VariablePool) GetString(name string) string {
	value := pool.Get(name)
	return fmt.Sprintf("%v", value)
}

func (pool *VariablePool) GetInt(name string) int {
	value := pool.Get(name)
	convertedValue, isOk := value.(int)

	if !isOk {
		pool.errorHandler.Capture(
			&errors.VariableCast{
				Name:     name,
				CastType: "int",
				RawValue: value,
			})
	}

	return convertedValue
}

func (pool *VariablePool) GetBool(name string) bool {
	value := pool.Get(name)
	convertedValue, isOk := value.(bool)

	if !isOk {
		pool.errorHandler.Capture(
			&errors.VariableCast{
				Name:     name,
				CastType: "bool",
				RawValue: value,
			})
	}

	return convertedValue
}

func (pool *VariablePool) Get(name string) interface{} {
	value, isPresent := pool.variables[name]

	if !isPresent {
		pool.errorHandler.Capture(&errors.VariableNotFound{Name: name})
	}

	return value
}

func (pool *VariablePool) Delete(name string) {
	delete(pool.variables, name)
}
