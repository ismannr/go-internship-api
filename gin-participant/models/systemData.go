package models

import (
	"gorm.io/gorm"
	"time"
)

type SystemData struct {
	gorm.Model
	Email          string    `json:"email"`
	Password       string    `json:"-" json:"password"`
	LastLogin      time.Time `json:"last_login"`
	Role           Role      `json:"role"`
	Level          Level     `json:"level"`
	CurrentlyLogin bool      `json:"currently_login"`
}
