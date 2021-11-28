package shooter

type Status string

const (
	Running      Status = "RUNNING"
	ShuttingDown Status = "SHUTTING_DOWN"
	Completed    Status = "COMPLETED"
	Stopped      Status = "STOPPED"
	Error        Status = "ERROR"
)
