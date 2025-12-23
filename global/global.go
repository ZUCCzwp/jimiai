package global

import (
	"github.com/go-pay/gopay/wechat/v3"
)

type (
	TraceIDKey struct{}
)

// WxClient 全局微信支付客户端
var WxClient *wechat.ClientV3
