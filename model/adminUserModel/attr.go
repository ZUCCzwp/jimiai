package adminUserModel

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type AttrEnum map[int]string

func (t *AttrEnum) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t AttrEnum) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type UserAttr struct {
	gorm.Model `json:"-"`
	SortId     int           `json:"sort_id"`
	AttrType   string        `json:"attr_type" gorm:"type:varchar(20)"`         // 属性类型
	AttrKey    string        `json:"attr_key" gorm:"<-:false;type:varchar(20)"` // 属性key
	AttrName   string        `json:"attr_name" gorm:"type:varchar(20)"`         // 属性名称
	Color      string        `json:"color" gorm:"type:varchar(20)"`             // 颜色
	HaveRange  bool          `json:"have_range"`                                // 是否有范围
	Range      UserAttrRange `json:"range" gorm:"embedded"`                     // 范围
}

type UserAttrRange struct {
	RangeType int      `json:"range_type"`             // 范围类型 0=范围 1=枚举
	Min       int      `json:"min"`                    // 最小值
	Max       int      `json:"max"`                    // 最大值
	Enum      AttrEnum `json:"enum" gorm:"type:json;"` // 枚举值
}

type UserAttrResponse struct {
	Id        int           `json:"id"`
	SortId    int           `json:"sort_id"`
	AttrType  string        `json:"attr_type"`
	AttrName  string        `json:"attr_name"`
	Color     string        `json:"color"`
	HaveRange bool          `json:"have_range"`
	Range     UserAttrRange `json:"range"`
}

type UserAttrRequest struct {
	UserAttrResponse
}
