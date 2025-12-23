package inviteModel

import "gorm.io/gorm"

// 用户排行榜结构体
type UserRank struct {
	UserId    int64   `json:"user_id"`    // 用户ID
	Score     float64 `json:"score"`      // 分数(机遇币 + 注册时间权重)
	Coins     int64   `json:"coins"`      // 机遇币数量
	JoinTime  int64   `json:"join_time"`  // 注册时间
	Nickname  string  `json:"nickname"`   // 昵称
	InviteNum int64   `json:"invite_num"` // 邀请人数
	Avatar    string  `json:"avatar"`     // 头像
	Rank      int64   `json:"rank"`       // 排名
}

// 邀请记录
type InviteLog struct {
	gorm.Model
	UserID       int64 `json:"user_id"`
	InviteUserID int64 `json:"invite_user_id"`
	InviteType   int   `json:"invite_type"` //邀请类型：0 免费 1vip
	InviteCoins  int   `json:"invite_coins"`
}

type InviteLogResponse struct {
	Id                 int64   `json:"id"`
	UserID             int64   `json:"user_id"`
	UserName           string  `json:"user_name"`
	FirstChargeNum     float64 `json:"first_charge_num"`
	InviteUserID       int64   `json:"invite_user_id"`
	InviteType         int     `json:"invite_type"` //邀请类型：0 免费 1vip
	InviteTime         string  `json:"invite_time"` // 邀请时间
	InviteCoins        int     `json:"invite_coins"`
	InviteUserNickname string  `json:"invite_user_nickname"` // 邀请人昵称
	InviteUserAvatar   string  `json:"invite_user_avatar"`   // 邀请人头像
}
