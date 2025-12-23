package chargeApi

import (
	"jiyu/model/chargeModel"
	"jiyu/model/contextModel"
	"jiyu/service/chargeService"
	"jiyu/service/redeemCodeService"
	"jiyu/util/response"
	"net/http"
	"strconv"

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

// RedeemList 充值明细列表接口
func RedeemList(c *gin.Context) {
	// 解析分页参数
	pn := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pn, 10, 64)
	if err != nil || page < 1 {
		response.Error(c, -1, "参数错误：page必须是正整数", nil)
		return
	}

	ps := c.DefaultQuery("pageSize", "20")
	pageSize, err := strconv.ParseInt(ps, 10, 64)
	if err != nil || pageSize < 1 {
		response.Error(c, -1, "参数错误：pageSize必须是正整数", nil)
		return
	}

	// 解析筛选参数
	code := c.Query("code")
	
	statusStr := c.Query("status")
	var status int = -1 // -1表示不筛选状态
	if statusStr != "" {
		statusInt, err := strconv.ParseInt(statusStr, 10, 64)
		if err != nil {
			response.Error(c, -1, "参数错误：status必须是数字", nil)
			return
		}
		status = int(statusInt)
	}

	// 调用service层获取数据
	list, count, err := redeemCodeService.GetRedeemCodeList(int(page), int(pageSize), code, status)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}
