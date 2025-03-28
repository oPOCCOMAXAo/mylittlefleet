package models

import "time"

type ContainerEvent struct {
	ID     int64           `gorm:"-"`
	Name   string          `gorm:"-"`
	Time   time.Time       `gorm:"-"`
	Status ContainerStatus `gorm:"-"`
}
