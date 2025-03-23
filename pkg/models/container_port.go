package models

type ContainerPort struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement"`
	ContainerID   int64 `gorm:"column:container_id;not null"`
	ContainerPort int64 `gorm:"column:container_port;not null"`
	HostPort      int64 `gorm:"column:host_port;not null"`
	IsPublic      bool  `gorm:"column:is_public;not null;default:false"`
}

func (ContainerPort) TableName() string {
	return "container_ports"
}

type ContainePortUniqueKey struct {
	ContainerID   int64
	ContainerPort int64
	HostPort      int64
	IsPublic      bool
}

func (p *ContainerPort) UniqueKey() ContainePortUniqueKey {
	return ContainePortUniqueKey{
		ContainerID:   p.ContainerID,
		ContainerPort: p.ContainerPort,
		HostPort:      p.HostPort,
		IsPublic:      p.IsPublic,
	}
}

func (p *ContainerPort) Equal(o *ContainerPort) bool {
	return p.ContainerID == o.ContainerID &&
		p.ContainerPort == o.ContainerPort &&
		p.HostPort == o.HostPort &&
		p.IsPublic == o.IsPublic
}

func (p *ContainerPort) PrepareForUpdate(o *ContainerPort) {
	p.ID = o.ID
}
