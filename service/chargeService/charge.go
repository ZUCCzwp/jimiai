package chargeService

import (
	"context"
	"fmt"
	"io"
	"jiyu/config"
	"jiyu/global"
	"jiyu/model/chargeModel"
	contextModel "jiyu/model/contextModel"
	"jiyu/model/payModel"
	"jiyu/repo/chargeRepo"
	"jiyu/repo/payRepo"
	"jiyu/service/memberService"
	"jiyu/service/userService"
	"jiyu/util"
	alipayCli "jiyu/util/pay/alipay"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	wxV3 "github.com/go-pay/gopay/wechat/v3"
	"github.com/smartwalle/apple"
	"github.com/spf13/cast"
)

func List(ctx contextModel.Context) (chargeModel.ListChargeResponse, error) {
	chargeList, err := chargeRepo.ListChargeConfig()
	if err != nil {
		log.Printf("get charge config err:%+v", err)
		return chargeModel.ListChargeResponse{}, err
	}
	// 转化price为元
	chargeListResponse := make([]chargeModel.ChargeInfoResponse, 0)
	for _, charge := range chargeList {
		chargeListResponse = append(chargeListResponse, chargeModel.ChargeInfoResponse{
			Price:       charge.GetPrice(),
			OriginPrice: charge.GetOriginPrice(),
			Month:       charge.Month,
		})
	}

	payment := make([]config.PaymentResponse, 0)
	payment = append(payment, config.PaymentResponse{
		Enable: config.PaymentConfig.AlipayInfo.AlipayEnable,
		Type:   "alipay",
	})
	payment = append(payment, config.PaymentResponse{
		Enable: config.PaymentConfig.WechatConfig.WechatEnable,
		Type:   "wechat",
	})
	return chargeModel.ListChargeResponse{ChargeList: chargeListResponse, Payment: payment}, nil
}

func MemberList(ctx contextModel.Context) (chargeModel.MemberListResponse, error) {
	memberList, err := memberService.ListMemberConfig()
	if err != nil {
		log.Printf("get member config err:%+v", err)
		return chargeModel.MemberListResponse{}, err
	}
	return chargeModel.MemberListResponse{
		MemberList: memberList,
		LevelRules: config.MemberConfig.LevelRules,
	}, nil
}

func CreateOrder(ctx contextModel.Context, req chargeModel.CreateOrderRequest) (chargeModel.CreateOrderResponse, error) {
	var (
		price     string
		err       error
		orderStr  string
		orderNo   = util.GenerateOrderId()
		nonceStr  string
		timeStamp string
		sign      string
	)
	// 根据支付类型获取支付金额
	chargeInfo := getChargeInfo(req.PayType)

	if req.PayMethod == "alipay" {
		price = chargeInfo.GetPrice()
		alipayClient := alipayCli.NewAliPayClient(config.PaymentConfig.AlipayInfo)
		alipayClient.SetLocation(alipay.LocationShanghai).
			SetCharset(alipay.UTF8).
			SetSignType(alipay.RSA2).
			SetNotifyUrl(config.PaymentConfig.AlipayInfo.NotifyUrl)

		bm := make(gopay.BodyMap)
		bm.Set("subject", "机遇充值包").
			Set("out_trade_no", orderNo).
			Set("total_amount", price)

		ctxV1 := context.Background()
		orderStr, err = alipayClient.TradeAppPay(ctxV1, bm)
		if err != nil {
			log.Printf("client.TradeAppPay(%+v),error:%+v", bm, err)
			return chargeModel.CreateOrderResponse{}, err
		}
	} else if req.PayMethod == "wechat" {
		price = cast.ToString(chargeInfo.GetWechatPrice())
		// 初始化 BodyMap
		expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
		bm := make(gopay.BodyMap)
		bm.Set("appid", config.PaymentConfig.WechatConfig.WechatAppID).
			Set("mchid", config.PaymentConfig.WechatConfig.WechatMchID).
			Set("description", "机遇充值包").
			Set("out_trade_no", cast.ToString(orderNo)).
			Set("time_expire", expire).
			Set("notify_url", config.PaymentConfig.WechatConfig.WechatNotifyUrl).
			SetBodyMap("amount", func(bm gopay.BodyMap) {
				bm.Set("total", cast.ToInt(price)).Set("currency", "CNY")
			})
		ctxV1 := context.Background()
		wxRsp, err := global.WxClient.V3TransactionApp(ctxV1, bm)
		if err != nil {
			log.Printf("wechatClient.V3TransactionApp(%+v),error:%+v", bm, err)
			return chargeModel.CreateOrderResponse{}, err
		}

		log.Printf("wxRsp:%s", wxRsp.Error)
		orderStr = wxRsp.Response.PrepayId

		apppayprams, err := global.WxClient.PaySignOfApp(config.PaymentConfig.WechatConfig.WechatAppID, orderStr)
		if err != nil {
			return chargeModel.CreateOrderResponse{}, err
		}

		nonceStr = apppayprams.Noncestr
		timeStamp = apppayprams.Timestamp
		sign = apppayprams.Sign
	}

	err = payRepo.SavePaymentLog(payModel.PaymentLog{
		Uid:         cast.ToInt(ctx.User.ID),
		OrderNo:     cast.ToString(orderNo),
		Rmb:         cast.ToFloat64(price),
		OrderStatus: payModel.PaymentStatusPending,
		PaymentType: req.PayMethod,
		PaymentEnv:  payModel.PaymentEnvApp,
		FreeVipDays: cast.ToInt(chargeInfo.Month * 30),
	})
	if err != nil {
		log.Printf("save payment log err:%+v", err)
		return chargeModel.CreateOrderResponse{}, err
	}

	return chargeModel.CreateOrderResponse{
		OrderStr:     orderStr,
		AppId:        config.PaymentConfig.WechatConfig.WechatAppID,
		PartnerId:    config.PaymentConfig.WechatConfig.WechatMchID,
		PackageValue: "Sign=WXPay",
		NonceStr:     nonceStr,
		TimeStamp:    timeStamp,
		Sign:         sign,
	}, nil
}

// 获取充值信息
// @author jimi
// @date 2025-01-04
// @param payType 支付类型
// @return config.ChargeInfo 充值信息
func getChargeInfo(payType int) chargeModel.ChargeInfo {
	chargeConfig, err := chargeRepo.GetChargeConfigByMonth(payType)
	if err != nil {
		log.Printf("get charge config err:%+v", err)
		return chargeModel.ChargeInfo{}
	}
	return chargeConfig
}

// 充值回调
func Notify(c *gin.Context) error {
	// 获取公钥
	aliPayPublicKey := config.PaymentConfig.AlipayInfo.AlipayPublicKey
	fmt.Println(888)
	fmt.Println(aliPayPublicKey)
	fmt.Println(999)
	bm, err := alipay.ParseNotifyToBodyMap(c.Request)
	if err != nil {
		log.Printf("parse alipay notify body err:%+v", err)
		return err
	}
	log.Printf("alipay notify body:%+v", bm)

	// 验签
	ok, err := alipay.VerifySign(aliPayPublicKey, bm)
	if err != nil {
		log.Printf("验签:%+v", err)
		return err
	}
	log.Printf("支付宝验签是否通过:%+v", ok)

	// 处理订单信息
	err = Paid(c, bm)
	if err != nil {
		log.Printf("处理订单信息失败:%+v", err)
		return err
	}

	return nil
}

func AppleNotify(c *gin.Context) error {
	var data, _ = io.ReadAll(c.Request.Body)
	var notification, err = apple.DecodeNotification([]byte(data))
	log.Printf("verify receipt:%+v", notification)

	if notification.Data != nil {
		// 处理订单信息
		err = ApplePaid(c, notification.Data.Transaction)
		if err != nil {
			log.Printf("处理订单信息失败:%+v", err)
			return err
		}
	}

	return nil
}

func WechatNotify(c *gin.Context) error {
	// 解析参数
	result, err := wxV3.V3ParseNotify(c.Request)
	if err != nil {
		log.Printf("parse wechat notify err:%+v", err)
		return err
	}
	log.Printf("wechat result:%+v", result)
	// certMap := client.WxPublicKeyMap()
	// bm := make(gopay.BodyMap)
	// bm.Set("appid", config.PaymentConfig.WechatConfig.WechatAppID).
	// 	Set("mchid", config.PaymentConfig.WechatConfig.WechatMchID).
	// 	Set("description", "机遇充值包").
	// 	Set("out_trade_no", cast.ToString(orderNo)).
	// 	Set("time_expire", expire).
	// 	Set("notify_url", config.PaymentConfig.WechatConfig.WechatNotifyUrl).
	// 	SetBodyMap("amount", func(bm gopay.BodyMap) {
	// 		bm.Set("total", cast.ToInt(price)).Set("currency", "CNY")
	// 	})
	// ctxV1 := context.Background()
	// wxRsp, err := global.WxClient.V3TransactionApp(ctxV1, bm)
	// client, err := wechat.NewClientV3(config.PaymentConfig.WechatConfig.WechatMchID, config.PaymentConfig.WechatConfig.WechatSerialNo, config.PaymentConfig.WechatConfig.WechatAPIv3Key, config.PaymentConfig.WechatConfig.WechatPrivateKey)
	// if err != nil {
	// 	xlog.Error(err)
	// 	return err
	// }

	// certMap := client.Lo
	// err = wxClient.LoadWechatPublicKeyByPath(publicKeyPath)
	// publicKey, err := loadPublicKey("path/to/wechatpay_public_key.pem")
	// pkMap := global.WxClient.WxPublicKeyMap()
	// publicKey := pkMap[result.SignInfo.HeaderSerial]
	// publicKey := global.WxClient.WxPublicKey()
	// err = wxV3.V3VerifySignByPK(result.SignInfo.HeaderTimestamp, result.SignInfo.HeaderNonce, result.SignInfo.SignBody, result.SignInfo.HeaderSignature, publicKey)
	// if err != nil {
	// 	log.Printf("verify wechat sign err:%+v", err)
	// 	return err
	// }
	// log.Printf("微信验签是否通过:%+v", true)

	text, err := result.DecryptPayCipherText(config.PaymentConfig.WechatConfig.WechatAPIv3Key)
	if err != nil {
		log.Printf("decrypt wechat cipher text err:%+v", err)
	}
	log.Printf("微信解密结果:%+v", text)

	// 处理订单信息
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", text.OutTradeNo)
	bm.Set("trade_no", text.TransactionId)
	err = Paid(c, bm)
	if err != nil {
		log.Printf("处理订单信息失败:%+v", err)
		return err
	}

	return nil
}

func Paid(ctx *gin.Context, bm gopay.BodyMap) error {
	// 获取订单信息
	orderNo := bm.GetString("out_trade_no")
	paymentLog, err := payRepo.FindByOrderNo(orderNo)
	if err != nil {
		log.Printf("find payment log by order no err:%+v", err)
		return err
	}
	// 开启事务
	tx := global.DB.Begin()
	// 更新订单
	err = payRepo.UpdatePaymentLog(tx, payModel.PaymentLog{
		OrderNo:     orderNo,
		OrderStatus: payModel.PaymentStatusSuccess,
		TradeNo:     bm.GetString("trade_no"),
	})
	if err != nil {
		log.Printf("update payment log err:%+v", err)
		tx.Rollback()
		return err
	}
	// 更新用户为vip并且增加会员的到期时间
	err = userService.UpdateUserVip(tx, paymentLog.Uid, paymentLog.FreeVipDays)
	if err != nil {
		log.Printf("update user vip err:%+v", err)
		tx.Rollback()
		return err
	}
	// 提交事务
	tx.Commit()

	return nil
}

func ApplePaid(ctx *gin.Context, transaction *apple.Transaction) error {
	// 获取订单信息
	orderNo := transaction.TransactionId

	paymentLog, err := payRepo.FindByOrderNo(orderNo)
	if err != nil {
		log.Printf("find payment log by order no err:%+v", err)
		return err
	}
	// 开启事务
	tx := global.DB.Begin()
	// 更新订单
	err = payRepo.UpdatePaymentLog(tx, payModel.PaymentLog{
		OrderNo:     orderNo,
		OrderStatus: payModel.PaymentStatusSuccess,
		TradeNo:     transaction.OriginalTransactionId,
	})
	if err != nil {
		log.Printf("update payment log err:%+v", err)
		tx.Rollback()
		return err
	}
	// 更新用户为vip并且增加会员的到期时间
	err = userService.UpdateUserVip(tx, paymentLog.Uid, paymentLog.FreeVipDays)
	if err != nil {
		log.Printf("update user vip err:%+v", err)
		tx.Rollback()
		return err
	}
	// 提交事务
	tx.Commit()

	return nil
}

func AdminChargeList(ctx contextModel.AdminContext, page, limit int) (result []chargeModel.ChargeInfo, count int, err error) {
	chargeList, err := chargeRepo.ListPageChargeConfig(page, limit)
	if err != nil {
		log.Printf("get charge config err:%+v", err)
		return nil, 0, err
	}

	count, err = chargeRepo.CountChargeConfig()
	if err != nil {
		return nil, 0, err
	}

	return chargeList, count, nil
}

func UpdateChargeConfig(ctx contextModel.AdminContext, charge chargeModel.ChargeInfo) error {
	return chargeRepo.UpdateChargeConfig(charge)
}
