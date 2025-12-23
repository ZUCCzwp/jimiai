package adminUserApi

import (
	"fmt"
	"jiyu/model/adminUserModel"
	"jiyu/model/contextModel"
	"jiyu/service/adminUserService"
	"jiyu/util/response"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func Ping(c *gin.Context) {
	response.Success(c, "pong", nil)
}

func Login(c *gin.Context) {
	var data adminUserModel.LoginRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "请填写用户名密码", nil)
	}

	jwt, err := adminUserService.Login(data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "登陆成功", gin.H{
		"token": jwt,
	})
}

func ChangePwd(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.ChangePwdRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.ChangePwd(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "修改成功", nil)
}

func Info(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	role := "管理员"

	if ctx.User.Role == 1 {
		role = "超级管理员"
	}

	result := adminUserModel.InfoResponse{
		Roles:        []string{role},
		Introduction: ctx.User.Introduction,
		Avatar:       ctx.User.Avatar,
		Name:         ctx.User.Name,
	}

	response.Success(c, "ok", result)
}

func Logout(c *gin.Context) {
	response.Success(c, "ok", nil)
}

func AdminUserList(c *gin.Context) {
	list, err := adminUserService.List()
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", list)
}

func AddAdminUser(c *gin.Context) {
	var data adminUserModel.AddUserRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err, data)
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := adminUserService.AddUser(data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "添加成功", nil)
}

func ResetAdminUserPassword(c *gin.Context) {
	var data adminUserModel.ResetPwdRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := adminUserService.ResetPwd(data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "重置成功", nil)
}

func DelAdminUser(c *gin.Context) {
	var data adminUserModel.UserIDRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := adminUserService.DeleteAdmin(data.Uid)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "删除成功", nil)
}

// ResetAdminUserAuth 重置管理员用户的权限
func ResetAdminUserAuth(c *gin.Context) {
	var data adminUserModel.ResetAuthorityRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err := adminUserService.ResetAuthority(data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "重置成功", nil)
}

func UserList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	pn := c.Query("page")
	page, err := strconv.ParseInt(pn, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	ps := c.Query("limit")
	limit, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	phone := c.Query("phone")
	nickname := c.Query("nickname")
	uid := c.Query("uid")
	list, count, err := adminUserService.UserList(ctx, int(page), int(limit), cast.ToInt(uid), phone, nickname)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}

func UserByUid(c *gin.Context) {
	u := c.Param("uid")
	uid, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	list, count, err := adminUserService.UserByUid(int(uid))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}

func Ban(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.UserBanRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.Ban(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func Delete(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.UserIDRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.Delete(ctx, data.Uid)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func Detail(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	id := c.Param("uid")
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	detail, err := adminUserService.Detail(ctx, int(uid))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", detail)
}

func Update(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.UpdateUserRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println(err)
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.Update(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func GetUserAttrs(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	attrs, err := adminUserService.GetAttrs(ctx)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}
	response.Success(c, "ok", attrs)
}

func GetUserAttrByKey(c *gin.Context) {
	key := c.Param("key")
	t := c.Param("type")

	attr, err := adminUserService.GetAttrByTypeAndKey(t, key)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", attr)
}

func UpdateUserAttrs(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.UserAttrRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.UpdateAttr(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func GetUserTags(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	pn := c.Query("page")
	page, err := strconv.ParseInt(pn, 10, 64)

	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	ps := c.Query("limit")
	limit, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	tags, count, err := adminUserService.GetUserTags(ctx, int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  tags,
		"total": count,
	})
}

func UpdateUserTags(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.UserTagRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.UpdateUserTag(ctx, data.Id, data.SortId, data.Title)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func CreateUserTags(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	var data adminUserModel.UserTagRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.CreateUserTag(ctx, data.SortId, data.Title, data.GroupName)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func DeleteUserTag(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	id := c.Query("id")
	tid, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = adminUserService.DeleteUserTag(ctx, int(tid))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}
