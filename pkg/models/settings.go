package models

type SettingsKey string

type Settings struct {
	Key   SettingsKey `gorm:"column:key;primaryKey;size:64;not null"`
	Value string      `gorm:"column:value;size:1024;not null"`
}

func (Settings) TableName() string {
	return "settings"
}
