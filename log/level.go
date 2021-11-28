package log

type Level struct {
	Order int
	Name  string
}

func (l Level) Above(refLevel Level) bool {
	return l.Order <= refLevel.Order
}

var (
	Error Level = Level{
		Order: 0,
		Name:  "ERROR",
	}
	Warning Level = Level{
		Order: 1,
		Name:  "WARNING",
	}
	Info Level = Level{
		Order: 2,
		Name:  "INFO",
	}
	Debug Level = Level{
		Order: 3,
		Name:  "DEBUG",
	}
)
