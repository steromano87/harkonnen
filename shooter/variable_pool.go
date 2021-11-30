package shooter

import (
	"fmt"
	"harkonnen/log"
)

type VariablePool struct {
	variables    map[string]interface{}
	logCollector *log.Collector
}

func NewVariablePool(logCollector *log.Collector) *VariablePool {
	pool := new(VariablePool)
	pool.variables = make(map[string]interface{})
	pool.logCollector = logCollector

	return pool
}

func (pool *VariablePool) Set(name string, value interface{}) {
	pool.variables[name] = value
	pool.logCollector.Debug(fmt.Sprintf("Set variable \"%s\" with value %s", name, value))
}

func (pool *VariablePool) GetString(name string) (string, error) {
	value, err := pool.Get(name)
	return fmt.Sprintf("%v", value), err
}

func (pool *VariablePool) GetInt(name string) (int, error) {
	value, err := pool.Get(name)
	if err != nil {
		return 0, err
	}

	convertedValue, isOk := value.(int)
	if !isOk {
		err = ErrBadVariableCast{
			Name:     name,
			CastType: "int",
			RawValue: value,
		}
		pool.logCollector.Error(err.Error())
		return 0, err
	}

	return convertedValue, nil
}

func (pool *VariablePool) GetBool(name string) (bool, error) {
	value, err := pool.Get(name)
	if err != nil {
		return false, err
	}

	convertedValue, isOk := value.(bool)
	if !isOk {
		err = ErrBadVariableCast{
			Name:     name,
			CastType: "bool",
			RawValue: value,
		}
		pool.logCollector.Error(err.Error())
		return false, err
	}

	return convertedValue, nil
}

func (pool *VariablePool) Get(name string) (interface{}, error) {
	pool.logCollector.Debug(fmt.Sprintf("Requested variable \"%s\"", name))
	value, isPresent := pool.variables[name]

	if !isPresent {
		err := ErrVariableNotFound{Name: name}
		pool.logCollector.Error(err.Error())
		return nil, err
	}

	return value, nil
}

func (pool *VariablePool) Delete(name string) {
	delete(pool.variables, name)
	pool.logCollector.Debug(fmt.Sprintf("Deleted variable \"%s\"", name))
}
