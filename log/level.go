package log

type Level struct {
	Order int
	Name  string
}

func (l Level) Above(refLevel Level) bool {
	return l.Order <= refLevel.Order
}

var (
	ErrorLevel = Level{
		Order: 0,
		Name:  "ERROR",
	}
	WarningLevel = Level{
		Order: 1,
		Name:  "WARNING",
	}
	InfoLevel = Level{
		Order: 2,
		Name:  "INFO",
	}
	DebugLevel = Level{
		Order: 3,
		Name:  "DEBUG",
	}
)
