package serverModel

import (
	"gorm.io/gorm"
	"time"
)

type Config struct {
	gorm.Model
	LastClearTime time.Time `json:"last_clear_time" gorm:"type:datetime;"`
}

func NewConfig() *Config {
	return &Config{
		LastClearTime: time.Now(),
	}
}
