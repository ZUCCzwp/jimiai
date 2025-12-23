package adminSettingApi

import (
	"github.com/gin-gonic/gin"
	"jiyu/model/contextModel"
	"jiyu/model/settingModel"
	"jiyu/service/settingService"
	"jiyu/util/response"
)

func Get(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	setting, err := settingService.Get(ctx)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}
	response.Success(c, "获取成功", setting)
}

func Update(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data settingModel.Setting
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := settingService.Update(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "修改成功", nil)
}
