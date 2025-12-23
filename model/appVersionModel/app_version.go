package appVersionModel

import (
	"gorm.io/gorm"
)

type AppVersion struct {
	gorm.Model
	Version      string `json:"version" gorm:"type:varchar(255);uniqueIndex:idx_version;"`
	BuildVersion string `json:"build_version" gorm:"type:varchar(255);"`
	Platform     string `json:"platform" gorm:"type:varchar(255);"`    // 平台: ios, android
	Status       bool   `json:"status" gorm:"type:bool;default:false"` // 是否强制更新
	Url          string `json:"url" gorm:"type:varchar(255);"`
	Content      string `json:"content" gorm:"type:text;"`
}
