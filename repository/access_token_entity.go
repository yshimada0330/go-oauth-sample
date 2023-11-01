package repository

import (
	"time"
)

type AccessToken struct {
	Id        uint64    `gorm:"column:id;primaryKey;"`
	Data      string    `gorm:"column:data;"`
	Token     string    `gorm:"column:token;"`
	ExpiredAt int64     `gorm:"column:expired_at;"`
	ClientId  string    `gorm:"column:client_id;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}
