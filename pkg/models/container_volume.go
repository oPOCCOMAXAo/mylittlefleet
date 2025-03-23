package models

type ContainerVolume struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	ContainerID   int64  `gorm:"column:container_id;not null"`
	VolumeID      int64  `gorm:"column:volume_id;not null"`
	ContainerPath string `gorm:"column:container_path;not null"`
}

func (ContainerVolume) TableName() string {
	return "container_volumes"
}

type ContainerVolumeUniqueKey struct {
	ContainerID   int64
	ContainerPath string
}

func (s *ContainerVolume) UniqueKey() ContainerVolumeUniqueKey {
	return ContainerVolumeUniqueKey{
		ContainerID:   s.ContainerID,
		ContainerPath: s.ContainerPath,
	}
}

func (s *ContainerVolume) Equal(o *ContainerVolume) bool {
	return s.ContainerID == o.ContainerID &&
		s.VolumeID == o.VolumeID &&
		s.ContainerPath == o.ContainerPath
}

func (s *ContainerVolume) PrepareForUpdate(o *ContainerVolume) {
	s.ID = o.ID
	s.VolumeID = o.VolumeID
}
