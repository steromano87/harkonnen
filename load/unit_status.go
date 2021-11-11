package load

type UnitStatus string

const (
	Running      UnitStatus = "RUNNING"
	ShuttingDown UnitStatus = "SHUTTING_DOWN"
	Completed    UnitStatus = "COMPLETED"
	Stopped      UnitStatus = "STOPPED"
	Error        UnitStatus = "ERROR"
)
