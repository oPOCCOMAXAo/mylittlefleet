package structs

import "github.com/opoccomaxao/mylittlefleet/pkg/models"

type ContainersDiff struct {
	DockerCreate []*models.FullContainerInfo
	DockerUpdate []*models.FullContainerInfo
	DockerDelete []*models.FullContainerInfo
	DockerStart  []*models.FullContainerInfo
	DockerStop   []*models.FullContainerInfo

	StorageUpdate []*models.FullContainerInfo
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
	Storage []*models.FullContainerInfo // storage containers, in db.
	Runtime []*models.FullContainerInfo // runtime containers, in docker.
}

type DiffContainerParams struct {
	Storage *models.FullContainerInfo // storage container, in db.
	Runtime *models.FullContainerInfo // runtime container, in docker.
}

type ContainersOptions struct {
	OnlyCurrentInstallation bool   // optional.
	InternalID              int64  // optional. ID in db.
	DockerID                string // optional.
	DockerName              string // optional.
}
