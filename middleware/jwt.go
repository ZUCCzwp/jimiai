package middleware

import (
	"fmt"
	"jiyu/global/redis"
	"jiyu/util"
	"jiyu/util/response"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth 验证JWT
func JWTAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			log.Printf("middleware.JWTAuth token为空,token: %s\n", token)
			response.Error(c, response.LoginAgain, "", nil)
			c.Abort()
			return
		}

		// 检查token是否在黑名单中
		key := redis.KeyTokenBlacklist + fmt.Sprintf("::%s", token)
		exists, err := redis.RDB.Exists(redis.Ctx, key).Result()
		if err == nil && exists > 0 {
			log.Printf("middleware.JWTAuth token在黑名单中, token: %s\n", token)
			response.Error(c, response.LoginAgain, "", nil)
			c.Abort()
			return
		}

		//打印token
		jwt, err := util.ValidJWT(token)
		if err != nil {
			log.Printf("middleware.JWTAuth valid jwt error,err: %v\n", err.Error())
			if strings.Contains(err.Error(), "token is expired") {
				// 登陆过期
				response.Error(c, response.LoginAgain, "", nil)
			} else {
				// token不合法
				response.Error(c, response.LoginAgain, "", nil)
			}

			c.Abort()
			return
		}

		fmt.Println("jwt: ", jwt)
		c.Set("uid", jwt.Uid)
		c.Set("nickname", jwt.Nickname)

		c.Next()
	}
}
