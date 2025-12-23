package redeemCodeApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/redeemCodeModel"
	"jiyu/service/redeemCodeService"
	"jiyu/util/response"
	"log"

	"github.com/gin-gonic/gin"
)

// CreateRedeemCode 创建兑换码接口
func CreateRedeemCode(c *gin.Context) {
	var req redeemCodeModel.CreateRedeemCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	// 验证参数
	if req.Amount <= 0 {
		response.Error(c, response.INVALID_PARAMS, "充值额度必须大于0", nil)
		return
	}

	result, err := redeemCodeService.CreateRedeemCode(req)
	if err != nil {
		log.Printf("redeemCodeApi.CreateRedeemCode 创建兑换码失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "创建成功", result)
}

// RedeemCode 兑换兑换码接口
func RedeemCode(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	var req redeemCodeModel.RedeemCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	// 验证参数
	if req.Code == "" {
		response.Error(c, response.INVALID_PARAMS, "兑换码不能为空", nil)
		return
	}

	err := redeemCodeService.RedeemCode(ctx, req)
	if err != nil {
		log.Printf("redeemCodeApi.RedeemCode 兑换兑换码失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "兑换成功", nil)
}
