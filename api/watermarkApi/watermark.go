package watermarkApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/watermarkModel"
	"jiyu/service/watermarkService"
	"jiyu/util/response"
	"log"

	"github.com/gin-gonic/gin"
)

// RemoveWatermark 去水印接口
func RemoveWatermark(c *gin.Context) {
	var req watermarkModel.ChatCompletionRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("watermarkApi.RemoveWatermark 参数错误, error: %v", err)
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	// 从 context 获取用户信息
	ctx := c.MustGet("context").(contextModel.Context)
	user := ctx.User

	// 从用户表获取 API Key
	apiKey := user.ApiKey
	if apiKey == "" {
		log.Printf("watermarkApi.RemoveWatermark 用户 API Key 未配置, userId: %d", user.ID)
		response.Error(c, response.ERROR, "API Key 未配置，请先设置您的 API Key", nil)
		return
	}

	// 调用 service 层
	result, err := watermarkService.RemoveWatermark(req, apiKey)
	if err != nil {
		log.Printf("watermarkApi.RemoveWatermark 调用服务失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	// 直接返回外部API的响应，保持格式一致
	response.Success(c, "ok", result)
}
