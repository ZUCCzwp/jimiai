package adminInviteApi

import (
	"jiyu/service/adminInviteService"
	"jiyu/service/inviteService"
	"jiyu/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InviteList(c *gin.Context) {
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

	inviteList, count, err := adminInviteService.AdminInviteList(int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  inviteList,
		"total": count,
	})
}

// InviteSuperior 查询上级
func InviteSuperior(c *gin.Context) {
	u := c.Param("uid")
	uid, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

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

	list, count, err := adminInviteService.InviteSuperior(int(uid), int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}
	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}

// InviteInferior 查询下级
func InviteInferior(c *gin.Context) {
	u := c.Param("uid")
	uid, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

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

	list, count, err := adminInviteService.InviteInferior(int(uid), int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}

func InviteRank(c *gin.Context) {
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

	list, count, err := inviteService.GetUserRankingByPage(int(page), int(limit))
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}
