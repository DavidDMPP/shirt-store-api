package domain

import (
	"time"
)

type MidtransConfig struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    ServerKey     string    `json:"server_key" gorm:"not null"`
    ClientKey     string    `json:"client_key" gorm:"not null"`
    Environment   string    `json:"environment" gorm:"not null;default:'sandbox'"`
    IsActive      bool      `json:"is_active" gorm:"default:true"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}