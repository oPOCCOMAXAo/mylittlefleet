package models

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/clients/docker"
	"github.com/samber/lo"
)

type ContainerStatus int

// CSStatic constant for static methods.
const CSStatic = ContainerStatus(0)

const (
	CSError ContainerStatus = iota
	CSStarting
	CSRunning
	CSStopping
	CSStopped
)

//nolint:gochecknoglobals
var containerStatusMap = map[ContainerStatus]string{
	CSError:    "Error",
	CSStarting: "Starting",
	CSRunning:  "Running",
	CSStopping: "Stopping",
	CSStopped:  "Stopped",
}

func (s ContainerStatus) String() string {
	return lo.CoalesceOrEmpty(containerStatusMap[s], "unknown")
}

//nolint:gochecknoglobals
var bsSubClassMap = map[ContainerStatus]string{
	CSError:    "danger",
	CSStarting: "success",
	CSRunning:  "success",
	CSStopping: "warning",
	CSStopped:  "danger",
}

func (s ContainerStatus) BSSubClass() string {
	return lo.CoalesceOrEmpty(bsSubClassMap[s], "secondary")
}

//nolint:gochecknoglobals
var dockerStateMap = map[string]ContainerStatus{
	docker.StateCreated:    CSStopped,
	docker.StateRunning:    CSRunning,
	docker.StatePaused:     CSStopped,
	docker.StateRestarting: CSStarting,
	docker.StateExited:     CSStopped,
	docker.StateRemoving:   CSStopped,
	docker.StateDead:       CSStopped,
}

func (ContainerStatus) FromDockerState(status string) ContainerStatus {
	return lo.CoalesceOrEmpty(dockerStateMap[status], CSError)
}
