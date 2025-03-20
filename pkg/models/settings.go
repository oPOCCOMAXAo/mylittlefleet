package models

type Settings struct {
	Key   string `gorm:"column:key;primaryKey;not null"`
	Value string `gorm:"column:value;not null"`
}

func (Settings) TableName() string {
	return "settings"
}
