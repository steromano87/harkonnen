package runtime

type VariablesHandler interface {
	Set(name string, value interface{})
	Delete(name string)
	Get(name string) interface{}
	GetString(name string) string
	GetInt(name string) int
	GetBool(name string) bool
}
