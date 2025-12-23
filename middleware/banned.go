package middleware

import (
	"jiyu/model/contextModel"
	"jiyu/util/response"

	"github.com/gin-gonic/gin"
)

// Banned 永久封禁
func Banned() func(c *gin.Context) {
	return func(c *gin.Context) {

		ctx := c.MustGet("context").(contextModel.Context)

		if ctx.User.BanLevel == 6 {
			response.Error(c, response.ERROR, "您已被封禁", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// Muted 被禁言
func Muted() func(c *gin.Context) {
	return func(c *gin.Context) {

		ctx := c.MustGet("context").(contextModel.Context)

		_, err := ctx.User.IsBanned()
		if err != nil {
			response.Error(c, response.ERROR, err.Error(), nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
