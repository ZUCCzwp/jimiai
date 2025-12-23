package adminUserModel

import "gorm.io/gorm"

type UserTag struct {
	gorm.Model
	SortId    int    `json:"sort_id"`
	Title     string `json:"title" gorm:"type:varchar(100)"`
	GroupName string `json:"group_name" gorm:"type:varchar(100)"`
}

type UserTagResponse struct {
	Id        int    `json:"id"`
	SortId    int    `json:"sort_id"`
	Title     string `json:"title"`
	GroupName string `json:"group_name"`
}

type UserTagRequest struct {
	UserTagResponse
}
