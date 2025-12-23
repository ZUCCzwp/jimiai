package inviteApi

import (
	"jiyu/model/contextModel"
	"jiyu/service/inviteService"
	"jiyu/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Ranking(c *gin.Context) {
	ranking, err := inviteService.GetUserRanking()
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", ranking)
}

func TopN(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	ranking, err := inviteService.GetTopNLeaderboard(n)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", ranking)
}

func UserRank(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)
	rank, err := inviteService.GetUserRankByUserId(int64(ctx.User.ID))
	if err != nil {
		response.Error(c, response.ERROR, "暂无该用户排名", nil)
		return
	}

	response.Success(c, "ok", rank)
}

func Record(c *gin.Context) {
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

	record, err := inviteService.GetInviteRecord(int64(ctx.User.ID), int(page), int(pageSize))
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", record)
}
