package models

type Volume struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	DockerName string `gorm:"column:docker_name;not null;default:''"`
	Internal   bool   `gorm:"column:internal;not null;default:false"`
}

func (Volume) TableName() string {
	return "volumes"
}
