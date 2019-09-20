package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

func Init(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Worker{},
		&UserWorker{},
		&Task{},
	).Error
}

type M map[string]interface{}

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
