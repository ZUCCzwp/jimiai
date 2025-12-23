package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
)

var (
	SUCCESS        = 20000
	ERROR          = -1
	LoginAgain     = 20001 // 重新登录
	INVALID_PARAMS = 20002 // 参数错误
)

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": SUCCESS,
		"msg":  msg,
		"data": data,
	})
}

func Error(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func NotifySuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gopay.SUCCESS)
}

func NotifyFail(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gopay.FAIL)
}
