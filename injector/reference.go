package injector

type Reference struct {
	Address string
	Port    uint16
	Weight  int
	Labels  []string
	Type
}
