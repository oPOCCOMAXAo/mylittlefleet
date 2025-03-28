package models

type Container struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement"`

	CreatedAt int64 `gorm:"column:created_at;not null;autoCreateTime"`

	Name string `gorm:"column:name;not null;index:idx_containers_name,unique"`

	// DockerID is used to identify the container in the docker.
	// Internal use only, don't expose to the user.
	DockerID string `gorm:"column:docker_id;not null;default:''"`

	// DockerName is used to identify the container in the docker.
	// Internal use only, don't expose to the user.
	DockerName string `gorm:"column:docker_name;not null;default:''"`

	Image string `gorm:"column:image;not null"`
	Tag   string `gorm:"column:tag;not null;default:''"`

	// Current status of the container.
	Paused bool `gorm:"column:paused;not null;default:true"`

	Deleted bool `gorm:"column:deleted;not null;default:false"`

	// Internal containers are not visible in the UI, but are used for internal purposes.
	Internal bool `gorm:"column:internal;not null;default:false"`
}

func (c *Container) TableName() string {
	return "containers"
}

type FullContainerInfo struct {
	Container *Container
	Volumes   []*VolumeDomain
	Ports     []*ContainerPort
	Envs      []*ContainerEnv
}
