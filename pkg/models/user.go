package models

type User struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
	Login     string `gorm:"column:login;unique;not null"`
	Password  string `gorm:"column:password;not null;size:128"` // bcrypt hashed.
}

func (User) TableName() string {
	return "users"
}
