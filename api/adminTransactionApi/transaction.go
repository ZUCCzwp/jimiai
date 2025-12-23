package adminTransactionApi

import (
	"github.com/gin-gonic/gin"
	"jiyu/util/response"
)

func List(c *gin.Context) {
	response.Success(c, "ok", gin.H{
		"items": make([]int, 0),
	})
}
