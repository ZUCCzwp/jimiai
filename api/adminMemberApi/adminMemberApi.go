package adminMemberApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/memberModel"
	"jiyu/service/memberService"
	"jiyu/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MemberConfigList(c *gin.Context) {
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

	memberList, count, err := memberService.AdminMemberConfigList(ctx, int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  memberList,
		"total": count,
	})
}

func MemberList(c *gin.Context) {
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

	memberList, count, err := memberService.AdminMemberList(ctx, int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  memberList,
		"total": count,
	})
}

func UpdateMemberConfig(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	member := memberModel.MemberConfig{}
	c.ShouldBindJSON(&member)

	err := memberService.UpdateMemberConfig(ctx, member)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func EditMember(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	member := memberModel.MemberEdit{}
	c.ShouldBindJSON(&member)

	err := memberService.EditMember(ctx, member)
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}
