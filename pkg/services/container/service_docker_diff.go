package container

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/updater"
)

// diffContainers compares runtime and storage containers and returns the difference.
// It does not modify the db directly, but changes structure from params.Storage.
// It returns the difference between runtime and storage containers.
func (s *Service) diffContainers(
	params structs.DiffContainersListParams,
) structs.ContainersDiff {
	var res structs.ContainersDiff

	rtcByID := make(map[int64]*structs.FullContainerInfo, len(params.Runtime))
	for _, container := range params.Runtime {
		rtcByID[container.Container.ID] = container
	}

	stcByID := make(map[int64]*structs.FullContainerInfo, len(params.Storage))
	for _, container := range params.Storage {
		stcByID[container.Container.ID] = container
	}

	for _, stc := range params.Storage {
		_, ok := rtcByID[stc.Container.ID]
		if !ok && !stc.Container.Deleted && !stc.Container.Paused {
			res.DockerCreate = append(res.DockerCreate, stc)
		}
	}

	for _, rtc := range params.Runtime {
		stc, ok := stcByID[rtc.Container.ID]
		if !ok || stc.Container.Deleted {
			res.DockerDelete = append(res.DockerDelete, rtc)

			continue
		}

		diff := s.diffRuntimeContainer(structs.DiffContainerParams{
			Storage: stc,
			Runtime: rtc,
		})

		res.Append(&diff)
	}

	return res
}

func (s *Service) diffRuntimeContainer(
	params structs.DiffContainerParams,
) structs.ContainersDiff {
	var res structs.ContainersDiff

	if params.Storage.Container.Paused {
		if !params.Runtime.Container.Paused {
			res.DockerStop = append(res.DockerStop, params.Runtime)
		}

		// we will apply updates on start.
		return res
	}

	if s.setStorageUpdate(params) {
		res.StorageUpdate = append(res.StorageUpdate, params.Storage)
	}

	if s.isRuntimeContainerChanged(params) {
		res.DockerUpdate = append(res.DockerUpdate, params.Storage)
	} else if params.Runtime.Container.Paused {
		res.DockerStart = append(res.DockerStart, params.Storage)
	}

	return res
}

// isRuntimeContainerChanged returns true if runtime container is different from storage container.
func (s *Service) isRuntimeContainerChanged(
	params structs.DiffContainerParams,
) bool {
	cmp := updater.NewComparer()

	updater.CompareValues(cmp, params.Storage.Container.Image, params.Runtime.Container.Image)
	updater.CompareValues(cmp, params.Storage.Container.Tag, params.Runtime.Container.Tag)
	updater.CompareValues(cmp, params.Storage.Container.DockerName, params.Runtime.Container.DockerName)

	s.isRuntimeVolumesChanged(cmp, params)
	s.isRuntimeEnvsChanged(cmp, params)
	s.isRuntimePortsChanged(cmp, params)

	return cmp.IsChanged()
}

func (s *Service) isRuntimeVolumesChanged(
	cmp *updater.Comparer,
	params structs.DiffContainerParams,
) {
	stvByPath := make(map[string]*structs.VolumeDomain, len(params.Storage.Volumes))
	for _, volume := range params.Storage.Volumes {
		stvByPath[volume.ContainerVolume.ContainerPath] = volume
	}

	rtvByPath := make(map[string]*structs.VolumeDomain, len(params.Runtime.Volumes))

	for _, rtv := range params.Runtime.Volumes {
		rtvByPath[rtv.ContainerVolume.ContainerPath] = rtv

		stv := stvByPath[rtv.ContainerVolume.ContainerPath]
		if stv == nil {
			cmp.SetChanged()

			continue
		}

		// if volume has no name, it is not updated after creation.
		if stv.Volume.DockerName != "" {
			updater.CompareValues(cmp, stv.Volume.DockerName, rtv.Volume.DockerName)
		}
	}

	for path := range stvByPath {
		rtv := rtvByPath[path]
		if rtv == nil {
			cmp.SetChanged()

			continue
		}
	}
}

func (s *Service) isRuntimeEnvsChanged(
	cmp *updater.Comparer,
	params structs.DiffContainerParams,
) {
	steByName := make(map[string]*models.ContainerEnv, len(params.Storage.Envs))
	for _, env := range params.Storage.Envs {
		steByName[env.Name] = env
	}

	rteByName := make(map[string]*models.ContainerEnv, len(params.Runtime.Envs))

	for _, rte := range params.Runtime.Envs {
		rteByName[rte.Name] = rte

		ste := steByName[rte.Name]
		if ste == nil {
			if !rte.IsDefault {
				cmp.SetChanged()
			}

			continue
		}

		updater.CompareValues(cmp, ste.Value, rte.Value)
	}

	for name := range steByName {
		rte := rteByName[name]
		if rte == nil {
			cmp.SetChanged()

			continue
		}
	}
}

func (s *Service) isRuntimePortsChanged(
	cmp *updater.Comparer,
	params structs.DiffContainerParams,
) {
	stpByUniqueID := make(map[models.ContainePortUniqueKey]*models.ContainerPort, len(params.Storage.Ports))
	for _, port := range params.Storage.Ports {
		stpByUniqueID[port.UniqueKey()] = port
	}

	rtpByUniqueID := make(map[models.ContainePortUniqueKey]*models.ContainerPort, len(params.Runtime.Ports))

	for _, rtp := range params.Runtime.Ports {
		rtpByUniqueID[rtp.UniqueKey()] = rtp

		stp := stpByUniqueID[rtp.UniqueKey()]
		if stp == nil {
			cmp.SetChanged()

			continue
		}
	}

	for key := range stpByUniqueID {
		rtp := rtpByUniqueID[key]
		if rtp == nil {
			cmp.SetChanged()

			continue
		}
	}
}
