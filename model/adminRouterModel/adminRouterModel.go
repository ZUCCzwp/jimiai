package adminRouterModel

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"jiyu/model"
)

type Children []RouterResponse

func (t *Children) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t Children) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type RouterBase struct {
	ParentId   int    `json:"parent_id" gorm:"type:int(11);not null;default:0;comment:'父路由id'"`
	Path       string `json:"path" gorm:"type:varchar(255);not null;comment:'路由路径'"`
	Name       string `json:"name" gorm:"type:varchar(255);not null;comment:'路由名称'"`
	Component  string `json:"component" gorm:"type:varchar(255);not null;comment:'路由组件'"`
	AlwaysShow bool   `json:"alwaysShow" gorm:"type:tinyint(1);not null;default:0;comment:'是否显示'"`
	Meta       Meta   `json:"meta" gorm:"embedded;"`
}

type Router struct {
	gorm.Model
	RouterBase
}

func (r Router) ConvertToResponse() RouterResponse {
	return RouterResponse{
		Id:         int(r.ID),
		RouterBase: r.RouterBase,
		Children:   make(Children, 0),
	}
}

type Meta struct {
	Title string        `json:"title" gorm:"type:varchar(255);not null;comment:'路由标题'"`
	Icon  string        `json:"icon" gorm:"type:varchar(255);not null;comment:'路由图标'"`
	Roles model.IntList `json:"roles" gorm:"type:json;not null;comment:'路由角色'"`
}

type RouterRequest struct {
	Id int `json:"id"`
	RouterBase
}

type RouterResponse struct {
	Id int `json:"id"`
	RouterBase
	Children Children `json:"children"`
}
