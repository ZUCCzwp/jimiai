package wechat

import (
	"jiyu/config"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/go-pay/xlog"
)

func NewWechatClient() *wechat.ClientV3 {
	client, err := wechat.NewClientV3(config.PaymentConfig.WechatConfig.WechatMchID, config.PaymentConfig.WechatConfig.WechatSerialNo, config.PaymentConfig.WechatConfig.WechatAPIv3Key, config.PaymentConfig.WechatConfig.WechatPrivateKey)
	if err != nil {
		xlog.Error(err)
		return nil
	}

	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = gopay.DebugOn
	return client
}
