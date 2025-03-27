package models

type ContainerEnv struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement"`
	ContainerID int64  `gorm:"column:container_id;not null;index:idx_envs_container"`
	Name        string `gorm:"column:name;not null;index:idx_envs_container"`
	Value       string `gorm:"column:value;not null"`

	// if true, this env is image default. It can't be removed.
	IsDefault    bool   `gorm:"column:is_default;not null;default:false"`
	DefaultValue string `gorm:"column:default_value;not null;default:''"`
}

func (ContainerEnv) TableName() string {
	return "container_envs"
}

type ContainerEnvUniqueKey struct {
	ContainerID int64
	Name        string
}

func (e *ContainerEnv) UniqueKey() ContainerEnvUniqueKey {
	return ContainerEnvUniqueKey{
		ContainerID: e.ContainerID,
		Name:        e.Name,
	}
}

func (e *ContainerEnv) Equal(other *ContainerEnv) bool {
	return e.ContainerID == other.ContainerID &&
		e.Name == other.Name &&
		e.Value == other.Value &&
		e.IsDefault == other.IsDefault &&
		e.DefaultValue == other.DefaultValue
}

func (e *ContainerEnv) PrepareForUpdate(other *ContainerEnv) {
	e.ID = other.ID
}
