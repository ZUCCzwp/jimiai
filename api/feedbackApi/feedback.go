package feedbackApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/globalModel"
	"jiyu/service/feedbackService"
	"jiyu/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func Add(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)
	var data globalModel.Feedback
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	data.UserId = int(ctx.User.ID)
	err = feedbackService.SaveFeedback(ctx, data)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func Read(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)
	feedbackId := c.Param("feedbackId")
	err := feedbackService.ReadFeedback(ctx, cast.ToInt(feedbackId))
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func List(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	pn := c.Query("page")
	page, err := strconv.ParseInt(pn, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	ps := c.Query("page_size")
	pageSize, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	result, count, err := feedbackService.GetList(ctx, int(page), int(pageSize))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  result,
		"total": count,
	})
}

func QAList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	pn := c.Query("page")
	page, err := strconv.ParseInt(pn, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	ps := c.Query("page_size")
	pageSize, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	result, count, err := feedbackService.GetQaList(ctx, int(page), int(pageSize))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  result,
		"total": count,
	})
}
