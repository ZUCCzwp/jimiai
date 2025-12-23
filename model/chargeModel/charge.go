package chargeModel

import (
	"jiyu/config"
	"jiyu/model/memberModel"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type ChargeInfo struct {
	gorm.Model
	Price       int `json:"price" gorm:"type:int"`        // 价格
	OriginPrice int `json:"origin_price" gorm:"type:int"` // 原价
	Month       int `json:"month" gorm:"type:int"`        // 月数:1-1个月会员、3-3个月会员、12-年会员
}

func (c *ChargeInfo) GetWechatPrice() int64 {
	return int64(c.Price)
}

func (c *ChargeInfo) GetPrice() string {
	return decimal.NewFromInt(cast.ToInt64(c.Price)).Div(decimal.NewFromInt(100)).String()
}

func (c *ChargeInfo) GetOriginPrice() string {
	return decimal.NewFromInt(cast.ToInt64(c.OriginPrice)).Div(decimal.NewFromInt(100)).String()
}

type ChargeInfoResponse struct {
	Price       string `json:"price"`
	OriginPrice string `json:"origin_price"`
	Month       int    `json:"month"`
}

type ListChargeResponse struct {
	ChargeList []ChargeInfoResponse     `json:"list"`
	Payment    []config.PaymentResponse `json:"payment"`
}

type ListPageChargeResponse struct {
	ChargeList []ChargeInfo `json:"list"`
	Count      int          `json:"count"`
}

type MemberListResponse struct {
	MemberList []memberModel.MemberConfig `json:"list"`
	LevelRules string                     `json:"level_rules"`
}

type CreateOrderRequest struct {
	PayMethod string `json:"pay_method"` // 支付方式：alipay、wechat、apple
	PayType   int    `json:"pay_type"`   // 支付类型：1-1个月、3-3个月、12-12个月
}

type CreateOrderResponse struct {
	OrderStr     string `json:"order_str"`
	AppId        string `json:"appid"`
	PartnerId    string `json:"partner_id"`
	PackageValue string `json:"package_value"`
	NonceStr     string `json:"nonce_str"`
	TimeStamp    string `json:"timestamp"`
	Sign         string `json:"sign"`
}
