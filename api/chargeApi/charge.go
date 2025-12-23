package chargeApi

import (
	"jiyu/model/chargeModel"
	"jiyu/model/contextModel"
	"jiyu/service/chargeService"
	"jiyu/util/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

func List(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	info, err := chargeService.List(ctx)

	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", info)
}

func MemberList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	info, err := chargeService.MemberList(ctx)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", info)
}

func CreateOrder(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	var req chargeModel.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	info, err := chargeService.CreateOrder(ctx, req)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", info)
}

func Notify(c *gin.Context) {
	err := chargeService.Notify(c)
	if err != nil {
		response.NotifyFail(c)
		return
	}

	response.NotifySuccess(c)
}

func AppleNotify(c *gin.Context) {
	err := chargeService.AppleNotify(c)
	if err != nil {
		response.NotifyFail(c)
		return
	}
	response.NotifySuccess(c)
}

func WechatNotify(c *gin.Context) {
	err := chargeService.WechatNotify(c)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
}
