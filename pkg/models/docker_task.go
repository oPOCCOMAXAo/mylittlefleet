package models

type DockerTaskAction string

const (
	DTAStart  DockerTaskAction = "start"
	DTAStop   DockerTaskAction = "stop"
	DTACreate DockerTaskAction = "create"
	DTADelete DockerTaskAction = "delete"
)

type DockerTask struct {
	ID          int64            `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt   int64            `gorm:"column:created_at;autoCreateTime"`
	Finished    bool             `gorm:"column:finished;default:false"`
	ContainerID int64            `gorm:"column:container_id"`
	Action      DockerTaskAction `gorm:"column:action"`
}

func (DockerTask) TableName() string {
	return "docker_tasks"
}
