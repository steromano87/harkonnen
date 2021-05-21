package runtime

type VariablesHandler interface {
	Set(name string, value interface{})
	Delete(name string)
	Get(name string) (value interface{}, isPresent bool)
	GetString(name string) (value string, isPresent bool)
	GetInt(name string) (value int, isPresent bool)
	GetBool(name string) (value bool, isPresent bool)
}
