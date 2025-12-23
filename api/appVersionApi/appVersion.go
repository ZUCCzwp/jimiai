package appVersionApi

import (
	"fmt"
	"jiyu/model/contextModel"
	versionmodel "jiyu/model/versionModel"
	"jiyu/service/appVersionService"
	versionservice "jiyu/service/versionService"
	"jiyu/util/response"

	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	platform := c.Query("platform")
	if platform == "" {
		response.Error(c, response.ERROR, "platform is required", nil)
		return
	}

	buildVersion := c.Query("build_version")

	version, err := appVersionService.GetVersion(platform, buildVersion)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", version)
}

func UpdateVersion(c *gin.Context) {
	var data versionmodel.AppVersion

	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("err:", err)
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}
	fmt.Println("data:", data)
	ctxValue := c.MustGet("context")
	ctx, ok := ctxValue.(contextModel.Context)
	if !ok {
		errMsg := "从上下文获取数据失败"
		fmt.Println(errMsg)
		response.Error(c, response.ERROR, errMsg, nil)
		return
	}
	err = versionservice.UpdateVersion(ctx, data)
	if err != nil {
		response.Error(c, response.ERROR, "更新昵称失败", nil)
	}
	response.Success(c, "ok", "ok")
}

// func UpadteNickName(c *gin.Context) {
// 	var data userModel.UserNickName

// 	err := c.ShouldBindJSON(&data)
// 	if err != nil {
// 		fmt.Println("err:", err)
// 		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
// 		return
// 	}

// 	ctx := c.MustGet("context").(contextModel.Context)

// 	err = userService.UpdateNickname(ctx, data)
// 	if err != nil {
// 		response.Error(c, response.ERROR, "更新昵称失败", nil)
// 	}

// 	response.Success(c, "ok", ctx.User.ID)
// }
