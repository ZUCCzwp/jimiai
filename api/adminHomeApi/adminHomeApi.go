package adminHomeApi

import (
	"jiyu/repo/payRepo"
	"jiyu/repo/userRepo"
	"jiyu/util/response"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	// 7日新增用户数
	count, err := userRepo.FindUserByIntervalDay(7)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	// 今日订单金额
	sum, err := payRepo.FindTodayRMB()
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	// 订单列表
	list, err := payRepo.ListForAdmin(1, 10, 0, -1, "", "")
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"count": count,
		"sum":   sum,
		"list":  list,
	})
}
