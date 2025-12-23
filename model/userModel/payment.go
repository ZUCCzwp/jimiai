package userModel

import (
	"database/sql/driver"
	"encoding/json"

	uuid "github.com/nu7hatch/gouuid"
)

type PaymentAccount struct {
	UUID           string `json:"uuid" gorm:"type:varchar(255);"`            // uuid
	PaymentType    string `json:"payment_type" gorm:"type:varchar(255);"`    // 支付类型 alipay wechat
	RealName       string `json:"real_name" gorm:"type:varchar(255);"`       // 真实姓名
	PaymentAccount string `json:"payment_account" gorm:"type:varchar(255);"` // 支付宝账号
}

type PaymentAccountList []PaymentAccount

func NewPaymentAccount(paymentType, realName, account string) PaymentAccount {
	u, _ := uuid.NewV4()
	return PaymentAccount{
		UUID:           u.String(),
		PaymentType:    paymentType,
		RealName:       realName,
		PaymentAccount: account,
	}
}

func (t *PaymentAccountList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t PaymentAccountList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type PaymentInfoRequest struct {
	PaymentType    string `json:"payment_type" gorm:"type:varchar(255);"`    // 支付类型 alipay wechat
	RealName       string `json:"real_name" gorm:"type:varchar(255);"`       // 真实姓名
	PaymentAccount string `json:"payment_account" gorm:"type:varchar(255);"` // 支付宝账号
}
