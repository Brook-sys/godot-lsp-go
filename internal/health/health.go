package health

type State string

const (
	Disconnected State = "disconnected"
	Connecting   State = "connecting"
	Connected    State = "connected"
	Reconnecting State = "reconnecting"
	Launching    State = "launching"
	ShuttingDown State = "shutting_down"
)
