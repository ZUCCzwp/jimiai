package userModel

import "gorm.io/gorm"

type HotUserResponse struct {
	ID       int    `json:"id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type HotUser struct {
	gorm.Model
	Uid           int    `json:"uid" gorm:"index:idx_uid"`
	Phone         string `json:"phone"`
	Shape         int    `json:"shape"`
	Nickname      string `json:"nickname"`
	AvatarUrl     string `json:"avatar_url"`
	Sort          int    `json:"sort"` // 排名
	HotValueCount int    `json:"hot_value_count"`
	HotConstCount int    `json:"hot_const_count"`
	Display       bool   `json:"display"`
}
