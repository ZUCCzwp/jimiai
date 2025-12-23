package payApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/payModel"
	"jiyu/model/userModel"
	"jiyu/repo/settingRepo"
	"jiyu/service/payService"
	"jiyu/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListForAdmin(c *gin.Context) {
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

	uid := c.Query("uid")
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		uidInt = 0
	}

	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	orderStatus := c.Query("order_status")
	orderStatusInt, err := strconv.ParseInt(orderStatus, 10, 64)
	if err != nil {
		orderStatusInt = 0
	}

	list, count, err := payService.ListForAdmin(ctx, int(page), int(limit), int(uidInt), int(orderStatusInt), startTime, endTime)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})

}

func WithdrawalListForAdmin(c *gin.Context) {
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

	uid := c.Query("uid")
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		uidInt = 0
	}

	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	status := c.Query("status")
	statusInt, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		statusInt = 0
	}

	list, count, err := payService.WithdrawalListForAdmin(ctx, int(page), int(limit), int(uidInt), int(statusInt), startTime, endTime)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})

}

func EditWithdrawal(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	id := c.Query("id")
	pid, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	s := c.Query("status")
	status, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	err = payService.EditWithdrawal(ctx, int(pid), int(status))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func AddAccount(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)
	var data userModel.PaymentInfoRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	err = payService.AddAccount(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func DeleteAccount(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)
	var data payModel.UUIDRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	err = payService.DeleteAccount(ctx, data.UUID)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func AccountList(c *gin.Context) {
	result := make([]interface{}, 0)
	response.Success(c, "ok", result)
}

func Withdrawal(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)
	var data payModel.WithdrawalRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	err = payService.Withdrawal(ctx, data)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "提交成功", nil)
}

func WithdrawalList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	list, err := payService.WithdrawalList(ctx)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", list)
}

func Setting(c *gin.Context) {
	setting, err := settingRepo.Find()
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"rate": setting.Withdraw.WithdrawRate,
		"min":  setting.Withdraw.MinWithdrawAmount,
	})
}
