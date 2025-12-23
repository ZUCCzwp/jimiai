package messageModel

import (
	"gorm.io/gorm"
	"jiyu/model"
)

type Message struct {
	gorm.Model
	MessageType    int              `gorm:"type:int;index:idx_message_type"` // 0:系统消息 1:漂流瓶点赞 2:漂流瓶评论 3:动态点赞 4:动态评论 5:关注
	Uid            int              `gorm:"type:int;not null;index:idx_uid"`
	Nickname       string           `gorm:"type:varchar(255);not null"`
	TargetUid      string           `gorm:"type:varchar(255);not null"`
	TargetNickname string           `gorm:"type:varchar(255);not null"`
	Content        string           `gorm:"type:varchar(255);not null"`
	Text           model.StringList `gorm:"type:json;not null"`
	Images         model.StringList `gorm:"type:json;not null"` // 相关图片, 如漂流瓶图片, 动态图片
	EntityId       int              `gorm:"type:int;not null"`  // 相关ID, 如漂流瓶ID, 动态ID
	IsRead         bool             `gorm:"type:bool;not null;default:false"`
}
