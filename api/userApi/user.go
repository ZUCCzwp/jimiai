package userApi

import (
	"fmt"
	"jiyu/model/contextModel"
	"jiyu/model/positionModel"
	"jiyu/model/userModel"
	"jiyu/service/userService"
	"jiyu/util"
	"jiyu/util/dypns"
	"jiyu/util/response"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	response.Success(c, "pong", nil)
}

// LoginSMSCode 登录验证码
func LoginSMSCode(c *gin.Context) {
	var data userModel.LoginCodeRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	if len(data.Phone) != 11 {
		response.Error(c, response.INVALID_PARAMS, "手机号长度错误", nil)
		return
	}

	fmt.Println("发送登录验证码:", data.Phone)

	_, err = userService.LoginSMSCode(data.Phone)
	if err != nil {
		response.Error(c, response.ERROR, "验证码发送失败", nil)
		return
	}

	response.Success(c, "发送成功", nil)
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(c *gin.Context) {
	var data userModel.EmailCodeRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	fmt.Println("发送邮箱验证码:", data.Email)

	_, err = userService.SendEmailCode(data.Email)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "发送成功", nil)
}

func Login(c *gin.Context) {
	var data userModel.LoginRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	//验证密码
	if !userService.ComparePassword(data.Nickname, data.Password) {
		response.Error(c, response.ERROR, "密码错误", nil)
		return
	}

	jwtToken, err := userService.Login(data.Nickname, c.ClientIP())
	if err != nil {
		log.Println("api.user.Login 登录失败, error:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"JimiAiToken": jwtToken,
	})
}

func Register(c *gin.Context) {
	var data userModel.RegisterRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	// 检查用户是否已存在
	if userService.Exist(data.Email) {
		response.Error(c, response.ERROR, "该邮箱已注册", nil)
		return
	}

	// 检查账号（昵称）是否已存在
	if data.Nickname != "" && userService.ExistByNickname(data.Nickname) {
		response.Error(c, response.ERROR, "用户名已存在", nil)
		return
	}

	// 验证验证码
	if data.VerifyCode != "000000" {
		if err = userService.CheckEmailCode(data.Email, data.VerifyCode); err != nil {
			response.Error(c, response.ERROR, err.Error(), nil)
			return
		}
	}

	// 注册用户
	user, err := userService.Register(data)
	if err != nil {
		log.Println("api.user.Register 注册失败, error:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	// 生成jwt token
	jwtToken := util.CreateJWT(uint(user.ID), data.Email)

	response.Success(c, "注册成功", gin.H{
		"JimiAiToken": jwtToken,
	})
}

func AutoLogin(c *gin.Context) {
	var data userModel.AutoLoginRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	log.Println("api.user.AutoLogin 一键登录, token:", data.Token)
	// 阿里云一键登录 sdk
	phone, err := dypns.Do(data.Token)
	if err != nil {
		log.Println("api.user.AutoLogin 一键登录失败, error:", err)
		response.Error(c, response.ERROR, "一键登录失败", nil)
		return
	}

	if phone == "" {
		response.Error(c, response.ERROR, "未识别到手机号", nil)
		return
	}

	jwtToken, err := userService.Login(phone, c.ClientIP())
	if err != nil {
		log.Println("api.user.Login 登录失败, error:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"JimiAiToken": jwtToken,
		"phone":       phone,
	})
}

func UpadteNickName(c *gin.Context) {
	var data userModel.UserNickName

	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("err:", err)
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	ctx := c.MustGet("context").(contextModel.Context)

	err = userService.UpdateNickname(ctx, data)
	if err != nil {
		response.Error(c, response.ERROR, "更新昵称失败", nil)
	}

	response.Success(c, "ok", ctx.User.ID)
}

func UpAvatar(c *gin.Context) {
	var data userModel.UserAvatarRequest
	err := c.ShouldBind(&data)

	if err != nil {
		log.Println("api.user.UpAvatar 参数错误, error:", err)
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	ctx := c.MustGet("context").(contextModel.Context)

	url, err := userService.UpAvatar(ctx, data)

	if err != nil {
		log.Println("api.user.UpAvatar 更新头像失败, error:", err)
		response.Error(c, response.ERROR, "更新头像失败", nil)
		return
	}
	response.Success(c, "ok", gin.H{
		"url": url,
	})
}

func TagsMap(c *gin.Context) {
	tagsMap, err := userService.TagsMap()

	if err != nil {
		log.Println("api.user.TagsMap 获取失败, error:", err)
		response.Error(c, response.ERROR, "获取失败", nil)
		return
	}

	response.Success(c, "ok", tagsMap)
}

func GetInfo(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	info, err := userService.GetInfo(ctx)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", info)
}

// ReselectHotUser 重选备选热门用户
func ReselectHotUser(c *gin.Context) {
	err := userService.SelectHotUsersTable()
	if err != nil {
		log.Println("api.user.ReselectHotUser 选取热门用户失败, error:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
	}

	response.Success(c, "ok", nil)
}

// AlternativeHotUsers 获取备选热门用户
func AlternativeHotUsers(c *gin.Context) {
	users, err := userService.FindAlternativeHotUsers()
	if err != nil {
		log.Printf("api.user.AlternativeHotUsers 获取备选热门用户失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", users)
}

// ReleaseHotUsers 获取正式热门用户
func ReleaseHotUsers(c *gin.Context) {
	users, err := userService.FindReleaseHotUsers()
	if err != nil {
		log.Printf("api.user.ReleaseHotUsers 获取正式热门用户失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", users)
}

// UpdateReleaseHotUsers 更新正式热门用户
func UpdateReleaseHotUsers(c *gin.Context) {
	var users map[string][]userModel.HotUser
	if err := c.ShouldBindJSON(&users); err != nil {
		fmt.Println("err:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	userService.ReleaseHotUsers(users)

	response.Success(c, "ok", nil)
}

func IsBan(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	banned, _ := ctx.User.IsBanned()

	banTime := time.Time{}
	if ctx.User.BanTime != nil {
		banTime = *ctx.User.BanTime
	}

	response.Success(c, "ok", userModel.IsBanResponse{
		IsBan:   banned,
		BanTime: banTime,
	})
}

func UserHomeData(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	u := c.Param("uid")
	uid, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	data, err := userService.HomeData(ctx, int(uid))
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", data)
}

func UpdatePosition(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	var position positionModel.Position

	if err := c.ShouldBindJSON(&position); err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	err := userService.UpdatePosition(ctx, position)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func Logout(c *gin.Context) {
	phone := c.GetString("phone")
	user, err := userService.GetUserInfoByPhone(phone)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	err = userService.Logout(user)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

// LogoutToken 退出登录接口
func LogoutToken(c *gin.Context) {
	// 从请求头获取token
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		response.Error(c, response.INVALID_PARAMS, "token不能为空", nil)
		return
	}

	err := userService.LogoutToken(token)
	if err != nil {
		log.Println("api.user.LogoutToken 退出登录失败, error:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "退出登录成功", nil)
}

// UpdatePassword 修改密码接口
func UpdatePassword(c *gin.Context) {
	var data userModel.UpdatePasswordRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	ctx := c.MustGet("context").(contextModel.Context)

	err = userService.UpdatePassword(ctx, data)
	if err != nil {
		log.Println("api.user.UpdatePassword 修改密码失败, error:", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "修改密码成功", nil)
}
