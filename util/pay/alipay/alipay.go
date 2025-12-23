package alipay

import (
	"jiyu/config"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/xlog"
)

func NewAliPayClient(config config.AliPayment) *alipay.Client {
	client, err := alipay.NewClient(config.AppId, config.PrivateKey, config.IsProd)
	if err != nil {
		xlog.Error(err)
		return nil
	}
	client.DebugSwitch = gopay.DebugOff
	if !config.IsProd {
		client.DebugSwitch = gopay.DebugOn
	}
	return client
}
