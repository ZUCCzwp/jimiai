package memberModel

import "gorm.io/gorm"

type MemberConfig struct {
	gorm.Model
	Type    int    `json:"type" gorm:"type:int;not null"`
	Name    string `json:"name" gorm:"type:varchar(255);not null"`
	Icon    string `json:"icon" gorm:"type:varchar(255);not null"`
	MinDays int    `json:"min_days" gorm:"type:int;not null"`
	MaxDays int    `json:"max_days" gorm:"type:int;not null"`
}

const (
	MemberTypeBronze  = 1
	MemberTypeSilver  = 2
	MemberTypeGolden  = 3
	MemberTypeDiamond = 4
)

type MemberListResponse struct {
	ID       int    `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	BuyTime  string `json:"buy_time"`
	VipTime  string `json:"vip_time"`
}
