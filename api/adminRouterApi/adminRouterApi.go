package adminRouterApi

import (
	"github.com/gin-gonic/gin"
	"jiyu/model/adminRouterModel"
	"jiyu/model/contextModel"
	"jiyu/service/adminRouterService"
	"jiyu/util/response"
)

func RouterList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	routers, err := adminRouterService.ListRouter(ctx.User.Role)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "获取成功", routers)
}

func SaveRouter(c *gin.Context) {
	var data adminRouterModel.RouterRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := adminRouterService.SaveRouter(data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "添加成功", nil)
}

func DeleteRouter(c *gin.Context) {
	var data struct {
		Id int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := adminRouterService.DelRouter(data.Id)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "删除成功", nil)
}
