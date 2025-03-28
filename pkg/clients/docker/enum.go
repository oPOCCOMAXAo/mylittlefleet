package docker

// Docker container states.
//
// See enum: https://docs.docker.com/reference/api/engine/version/v1.48/#tag/Container/operation/ContainerList
const (
	StateUnknown    = ""
	StateCreated    = "created"
	StateRunning    = "running"
	StatePaused     = "paused"
	StateRestarting = "restarting"
	StateExited     = "exited"
	StateRemoving   = "removing"
	StateDead       = "dead"
)
