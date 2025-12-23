package middleware

import (
	"jiyu/model/contextModel"
	"jiyu/model/userModel"
	"jiyu/repo/adminUserRepo"
	"jiyu/repo/userRepo"
	"jiyu/util/response"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ContextAdmin() func(c *gin.Context) {
	return func(c *gin.Context) {
		context := contextModel.AdminContext{}

		phone := c.GetString("phone")

		user, err := adminUserRepo.FindByUsername(phone)
		if err != nil && user.ID > 0 {
			log.Printf("middleware.Context 用户获取失败, phone: %s, error: %s\n", phone, err)
			response.Error(c, response.ERROR, "用户获取失败", nil)
			c.Abort()
			return
		}

		context.User = user

		c.Set("context", context)

		c.Next()
	}
}

func Context() func(c *gin.Context) {
	return func(c *gin.Context) {
		context := contextModel.Context{}

		nickname := c.GetString("nickname")

		user, err := userRepo.FindByUserName(nickname)
		if err != nil && user.ID > 0 {
			log.Printf("middleware.Context 用户获取失败, nickname: %s, error: %s\n", nickname, err)
			response.Error(c, response.ERROR, "用户获取失败", nil)
			c.Abort()
			return
		}

		context.User = user

		// 更新用户登陆信息
		updateUser := userModel.User{
			Model:           gorm.Model{ID: uint(user.ID)},
			LastLoginIp:     c.ClientIP(),
			LastRequestTime: time.Now(),
		}

		err = userRepo.Updates(&updateUser)
		if err != nil {
			log.Printf("middleware.Context 更新用户数据登陆信息失败, nickname: %s, error: %s\n", nickname, err)
			response.Error(c, response.ERROR, "更新用户数据登陆信息失败", nil)
			c.Abort()
			return
		}

		c.Set("context", context)

		c.Next()
	}
}
