package redeemCodeModel

import (
	"time"

	"gorm.io/gorm"
)

// RedeemCodeStatus 兑换码状态
const (
	RedeemCodeStatusUnused = 0 // 未使用
	RedeemCodeStatusUsed   = 1 // 已使用
)

// RedeemCode 兑换码表
type RedeemCode struct {
	gorm.Model
	Code       string     `json:"code" gorm:"type:varchar(50);not null;unique;comment:'兑换码'"`             // 兑换码
	Amount     int        `json:"amount" gorm:"type:int;not null;default:0;comment:'充值额度'"`               // 充值额度
	Status     int        `json:"status" gorm:"type:tinyint;not null;default:0;comment:'状态:0-未使用,1-已使用'"` // 状态
	UserId     int        `json:"user_id" gorm:"type:int;default:0;comment:'兑换用户ID'"`                     // 兑换用户ID
	RedeemedAt *time.Time `json:"redeemed_at" gorm:"type:datetime;default:null;comment:'兑换时间'"`           // 兑换时间
}

// CreateRedeemCodeRequest 创建兑换码请求
type CreateRedeemCodeRequest struct {
	Amount int `json:"amount" binding:"required"` // 充值额度
}

// RedeemCodeRequest 兑换兑换码请求
type RedeemCodeRequest struct {
	Code string `json:"code" binding:"required"` // 兑换码
}

// RedeemCodeResponse 兑换码响应
type RedeemCodeResponse struct {
	Code      string    `json:"code"`       // 兑换码
	Amount    int       `json:"amount"`     // 充值额度
	Status    int       `json:"status"`     // 状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// RedeemCodeListResponse 充值明细列表响应
type RedeemCodeListResponse struct {
	Code         string     `json:"code"`          // 兑换码
	Status       int        `json:"status"`        // 兑换状态 0-未使用 1-已使用
	UserNickname string     `json:"user_nickname"` // 兑换人昵称
	RedeemedAt   *time.Time `json:"redeemed_at"`   // 兑换时间
	Amount       int        `json:"amount"`        // 充值金额
}
