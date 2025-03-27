package container

import (
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/container/structs"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/updater"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/xslices"
)

// setStorageUpdate updates storage container struct with runtime data.
//
// It returns true if storage container was updated.
//
// Does not update data in the database.
func (s *Service) setStorageUpdate(
	params structs.DiffContainerParams,
) bool {
	upd := updater.New()

	updater.SetValue(upd, &params.Storage.Container.DockerID, params.Runtime.Container.DockerID)

	s.setStorageEnvsUpdate(upd, params)
	s.setStorageVolumesUpdate(upd, params)

	return upd.IsChanged()
}

func (s *Service) setStorageVolumesUpdate(
	upd *updater.Updater,
	params structs.DiffContainerParams,
) {
	rtvByPath := make(map[string]*structs.VolumeDomain, len(params.Runtime.Volumes))
	for _, volume := range params.Runtime.Volumes {
		rtvByPath[volume.ContainerVolume.ContainerPath] = volume
	}

	for i, vol := range params.Storage.Volumes {
		rtv := rtvByPath[vol.ContainerVolume.ContainerPath]
		if rtv == nil {
			params.Storage.Volumes[i] = nil

			upd.SetChanged()

			continue
		}

		updater.SetValue(upd, &vol.Volume.DockerName, rtv.Volume.DockerName)
	}

	xslices.RemoveZeroRef(&params.Storage.Volumes)
}

func (s *Service) setStorageEnvsUpdate(
	upd *updater.Updater,
	params structs.DiffContainerParams,
) {
	rteByName := make(map[string]*models.ContainerEnv, len(params.Runtime.Envs))
	for _, env := range params.Runtime.Envs {
		rteByName[env.Name] = env
	}

	for i, env := range params.Storage.Envs {
		rte := rteByName[env.Name]
		if rte == nil {
			params.Storage.Envs[i] = nil

			upd.SetChanged()

			continue
		}

		updater.SetValue(upd, &env.Value, rte.Value)
		updater.SetValue(upd, &env.IsDefault, rte.IsDefault)
		updater.SetValue(upd, &env.DefaultValue, rte.DefaultValue)
	}

	xslices.RemoveZeroRef(&params.Storage.Envs)
}
