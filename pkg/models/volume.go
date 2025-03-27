package models

type Volume struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	DockerName string `gorm:"column:docker_name;not null;default:''"`
	Internal   bool   `gorm:"column:internal;not null;default:false"`
}

func (Volume) TableName() string {
	return "volumes"
}

func (v *Volume) UniqueKey() int64 {
	return v.ID
}

func (v *Volume) Equal(other *Volume) bool {
	return v.DockerName == other.DockerName &&
		v.Internal == other.Internal
}

func (v *Volume) PrepareForUpdate(_ *Volume) {
}
