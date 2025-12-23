package adminChargeApi

import (
	"jiyu/model/chargeModel"
	"jiyu/model/contextModel"
	"jiyu/service/chargeService"
	"jiyu/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ChargeConfigList(c *gin.Context) {
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

	chargeList, count, err := chargeService.AdminChargeList(ctx, int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  chargeList,
		"total": count,
	})
}

func UpdateChargeConfig(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	charge := chargeModel.ChargeInfo{}
	c.ShouldBindJSON(&charge)

	err := chargeService.UpdateChargeConfig(ctx, charge)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}
