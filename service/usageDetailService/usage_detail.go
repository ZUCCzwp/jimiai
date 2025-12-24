package usageDetailService

import (
	"errors"
	"jiyu/model/contextModel"
	"jiyu/model/usageDetailModel"
	"jiyu/repo/usageDetailRepo"
	"log"
)

// CreateUsageDetail 创建使用明细
func CreateUsageDetail(ctx contextModel.Context, req usageDetailModel.CreateUsageDetailRequest) error {
	userId := int(ctx.User.ID)
	user := ctx.User

	// 将英文类型转换为中文类型（数据库存储使用中文）
	typeCN := usageDetailModel.ConvertTypeToCN(req.Type)

	// 计算实际花费
	var actualCost float64

	// 如果是消费类型，根据模型和用户VIP状态自动计算花费
	if req.Type == usageDetailModel.UsageDetailTypeConsumption && req.Model != "" {
		isVip := user.IsVip()
		actualCost = CalculateCost(req.Model, isVip)
		log.Printf("usageDetailService.CreateUsageDetail 计算花费: 模型=%s, VIP=%v, 花费=%.4f", req.Model, isVip, actualCost)
	} else {
		// 充值、退款等其他类型，使用请求中的花费金额
		actualCost = req.Cost
		if actualCost <= 0 {
			return errors.New("花费金额必须大于0")
		}
	}

	usageDetail := &usageDetailModel.UsageDetail{
		UserId:    userId,
		Type:      typeCN, // 存储中文类型
		ModelName: req.Model,
		Cost:      actualCost,
		Details:   req.Details,
	}

	err := usageDetailRepo.CreateUsageDetail(usageDetail)
	if err != nil {
		log.Printf("usageDetailService.CreateUsageDetail 创建使用明细失败, error: %v", err)
		return errors.New("创建使用明细失败")
	}

	return nil
}

// GetUsageDetailList 获取使用明细列表
func GetUsageDetailList(ctx contextModel.Context, req usageDetailModel.UsageDetailListRequest) ([]usageDetailModel.UsageDetailListResponse, int, error) {
	userId := int(ctx.User.ID)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 100
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	list, err := usageDetailRepo.GetUsageDetailList(userId, req.Page, req.PageSize, req.StartTime, req.EndTime, req.Type, req.Model)
	if err != nil {
		log.Printf("usageDetailService.GetUsageDetailList 获取使用明细列表失败, error: %v", err)
		return nil, 0, errors.New("获取使用明细列表失败")
	}

	count, err := usageDetailRepo.CountUsageDetailList(userId, req.StartTime, req.EndTime, req.Type, req.Model)
	if err != nil {
		log.Printf("usageDetailService.GetUsageDetailList 统计使用明细列表总数失败, error: %v", err)
		return nil, 0, errors.New("统计使用明细列表总数失败")
	}

	return list, count, nil
}
