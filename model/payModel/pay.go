package payModel

import (
	"time"

	"gorm.io/gorm"
)

type UUIDRequest struct {
	UUID string `json:"uuid"`
}

const (
	PaymentTypeAlipay = "alipay"
	PaymentTypeWechat = "wechat"
	PaymentTypeApple  = "apple"
)

const (
	PaymentEnvApp = "app"
	PaymentEnvWeb = "web"
	PaymentEnvH5  = "h5"
)

const (
	PaymentStatusPending = 0
	PaymentStatusSuccess = 1
	PaymentStatusFailed  = 2
)

// PaymentLog 支付日志
type PaymentLog struct {
	gorm.Model
	Uid         int     `json:"uid" gorm:"type:int;not null;index:idx_uid;"`
	Nickname    string  `json:"nickname" gorm:"type:varchar(255);not null;"`
	Rmb         float64 `json:"rmb" gorm:"default:0"`
	Diamond     int     `json:"diamond" gorm:"type:int;not null;"`
	FreeDiamond int     `json:"free_diamond" gorm:"type:int;not null;"`              // 赠送钻石
	OrderNo     string  `json:"order_no" gorm:"type:varchar(255);not null;"`         // 商户订单号
	TradeNo     string  `json:"trade_no" gorm:"type:varchar(255);not null;"`         // 第三方订单号
	PaymentType string  `json:"payment_type" gorm:"type:varchar(255);not null;"`     // 支付类型 alipay wechat
	PaymentEnv  string  `json:"payment_env" gorm:"type:varchar(255);not null;"`      // 支付环境 app web
	FreeJifen   int     `json:"free_jifen" gorm:"type:int;not null;"`                // 赠送积分
	FreeVipDays int     `json:"free_vip_days" gorm:"type:int;not null;"`             // 赠送vip天数
	FreeGift    int     `json:"free_gift" gorm:"type:int;not null;"`                 // 赠送礼物
	OrderStatus int     `json:"order_status" gorm:"type:int;not null;default:0"`     // 订单状态 0 未支付 1 已支付
	Remark      string  `json:"remark" gorm:"type:varchar(255);not null;default:''"` // 备注
}

type WithdrawalLog struct {
	gorm.Model
	Uid               int       `json:"uid" gorm:"type:int;not null;index:idx_uid;"`
	Nickname          string    `json:"nickname" gorm:"type:varchar(255);not null;"`
	Rmb               float64   `json:"rmb" gorm:"type:int;not null;"`
	Ticket            int       `json:"ticket" gorm:"type:int;not null;"`                      // 兑换映票
	WithdrawalType    string    `json:"withdrawal_type" gorm:"type:varchar(255);not null;"`    // 提现类型 alipay wechat
	WithdrawalName    string    `json:"withdrawal_name" gorm:"type:varchar(255);not null;"`    // 提现姓名
	WithdrawalAccount string    `json:"withdrawal_account" gorm:"type:varchar(255);not null;"` // 提现账号
	Status            int       `json:"status" gorm:"type:int;not null;default:0"`             // 提现状态 0 未处理 1 已处理 2 已拒绝
	HandleTime        time.Time `json:"handle_time" gorm:"type:datetime;not null;"`            // 处理时间
}

type WithdrawalRequest struct {
	AccountUUID string  `json:"account_uuid"`
	RMB         float64 `json:"rmb"`
}

type WithdrawalListResponse struct {
	Id      int       `json:"id"`
	RMB     float64   `json:"rmb"`
	Account string    `json:"account"`
	Time    time.Time `json:"time"`
	Status  string    `json:"status"`
}
