package structs

import "github.com/opoccomaxao/mylittlefleet/pkg/models"

type FullContainerInfo struct {
	Container *models.Container
	Volumes   []*VolumeDomain
	Ports     []*models.ContainerPort
	Envs      []*models.ContainerEnv
}

type ContainersDiff struct {
	DockerCreate []*FullContainerInfo
	DockerUpdate []*FullContainerInfo
	DockerDelete []*FullContainerInfo
	DockerStart  []*FullContainerInfo
	DockerStop   []*FullContainerInfo

	StorageUpdate []*FullContainerInfo
}

func (diff *ContainersDiff) Append(other *ContainersDiff) {
	diff.DockerCreate = append(diff.DockerCreate, other.DockerCreate...)
	diff.DockerUpdate = append(diff.DockerUpdate, other.DockerUpdate...)
	diff.DockerDelete = append(diff.DockerDelete, other.DockerDelete...)
	diff.DockerStart = append(diff.DockerStart, other.DockerStart...)
	diff.DockerStop = append(diff.DockerStop, other.DockerStop...)
	diff.StorageUpdate = append(diff.StorageUpdate, other.StorageUpdate...)
}

type DiffContainersListParams struct {
	Storage []*FullContainerInfo // storage containers, in db.
	Runtime []*FullContainerInfo // runtime containers, in docker.
}

type DiffContainerParams struct {
	Storage *FullContainerInfo // storage container, in db.
	Runtime *FullContainerInfo // runtime container, in docker.
}

type ContainersOptions struct {
	OnlyCurrentInstallation bool   // optional.
	InternalID              int64  // optional. ID in db.
	DockerID                string // optional.
	DockerName              string // optional.
}
