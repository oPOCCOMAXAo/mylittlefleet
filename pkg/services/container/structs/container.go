package structs

import "github.com/opoccomaxao/mylittlefleet/pkg/models"

type FullContainerInfo struct {
	Container *models.Container
	Volumes   []*models.ContainerVolume
	Ports     []*models.ContainerPort
	Envs      []*models.ContainerEnv
}
