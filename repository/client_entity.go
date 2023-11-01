package repository

import (
	"time"
)

type Client struct {
	Id        uint64    `gorm:"column:id;primaryKey;"`
	ClientId  string    `gorm:"column:client_id;"`
	Secret    string    `gorm:"column:secret;"`
	Domain    string    `gorm:"column:domain;"`
	Scope     string    `gorm:"column:scope;"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}
