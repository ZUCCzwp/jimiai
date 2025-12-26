package usageDetailApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/usageDetailModel"
	"jiyu/service/usageDetailService"
	"jiyu/util/response"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUsageDetail 创建使用明细接口
func CreateUsageDetail(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	var req usageDetailModel.CreateUsageDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	// 验证参数
	if req.Type == "" {
		response.Error(c, response.INVALID_PARAMS, "类型不能为空", nil)
		return
	}

	// 验证类型是否有效
	validTypes := map[string]bool{
		usageDetailModel.UsageDetailTypeConsumption: true,
		usageDetailModel.UsageDetailTypeRefund:      true,
		usageDetailModel.UsageDetailTypeCharge:      true,
	}
	if !validTypes[req.Type] {
		response.Error(c, response.INVALID_PARAMS, "类型无效，必须是: consumption(消费)、refund(退款)、recharge(充值)", nil)
		return
	}

	// 如果是消费类型，必须提供模型名称或任务ID
	if req.Type == usageDetailModel.UsageDetailTypeConsumption && req.Model == "" && req.TaskId == "" {
		response.Error(c, response.INVALID_PARAMS, "消费类型必须提供模型名称或任务ID", nil)
		return
	}

	err := usageDetailService.CreateUsageDetail(ctx, req)
	if err != nil {
		log.Printf("usageDetailApi.CreateUsageDetail 创建使用明细失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "创建成功", nil)
}

// GetUsageDetailList 获取使用明细列表接口
func GetUsageDetailList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	// 解析分页参数
	pn := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pn, 10, 64)
	if err != nil || page < 1 {
		response.Error(c, response.INVALID_PARAMS, "参数错误：page必须是正整数", nil)
		return
	}

	ps := c.DefaultQuery("pageSize", "100")
	pageSize, err := strconv.ParseInt(ps, 10, 64)
	if err != nil || pageSize < 1 {
		response.Error(c, response.INVALID_PARAMS, "参数错误：pageSize必须是正整数", nil)
		return
	}

	// 解析筛选参数
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	detailType := c.Query("type")
	modelName := c.Query("modelName")

	req := usageDetailModel.UsageDetailListRequest{
		Page:      int(page),
		PageSize:  int(pageSize),
		StartTime: startTime,
		EndTime:   endTime,
		Type:      detailType,
		Model:     modelName,
	}

	// 调用service层获取数据
	list, count, err := usageDetailService.GetUsageDetailList(ctx, req)
	if err != nil {
		log.Printf("usageDetailApi.GetUsageDetailList 获取使用明细列表失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}
