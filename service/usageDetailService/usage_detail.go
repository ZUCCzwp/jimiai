package usageDetailService

import (
	"errors"
	"fmt"
	"jiyu/model/contextModel"
	"jiyu/model/usageDetailModel"
	"jiyu/repo/taskRepo"
	"jiyu/repo/usageDetailRepo"
	"log"
)

// CreateUsageDetail 创建使用明细
func CreateUsageDetail(ctx contextModel.Context, req usageDetailModel.CreateUsageDetailRequest) error {
	userId := int(ctx.User.ID)
	user := ctx.User

	// 计算实际花费
	var actualCost float64
	var modelName string = req.Model

	// 如果提供了任务ID，从任务中获取模型信息
	if req.TaskId != "" {
		task, err := taskRepo.GetTaskByTaskId(req.TaskId)
		if err != nil {
			log.Printf("usageDetailService.CreateUsageDetail 获取任务信息失败, taskId=%s, error: %v", req.TaskId, err)
			return errors.New("任务不存在")
		}

		// 验证任务是否属于当前用户
		if task.UserId != userId {
			return errors.New("无权操作此任务")
		}

		// 从任务中获取模型名称（优先使用任务中的模型名称）
		if task.ModelName != "" {
			modelName = task.ModelName
		}
		// 如果任务中没有模型名称，使用请求中的模型名称（如果有的话，已在初始化时设置）
	}

	// 如果是消费类型或退款类型且提供了模型名称，根据模型和用户VIP状态自动计算花费
	// 注意：任务失败（退款类型+提供taskId）时，花费金额会根据模型自动计算，不使用请求中的cost字段
	if modelName != "" && (req.Type == usageDetailModel.UsageDetailTypeConsumption || req.Type == usageDetailModel.UsageDetailTypeRefund) {
		isVip := user.IsVip()
		actualCost = CalculateCost(modelName, isVip)
		log.Printf("usageDetailService.CreateUsageDetail 计算花费: 类型=%s, 任务ID=%s, 模型=%s, VIP=%v, 花费=%.4f", req.Type, req.TaskId, modelName, isVip, actualCost)
	} else {
		// 充值等其他类型，或没有提供模型名称的退款/消费类型，使用请求中的花费金额
		actualCost = req.Cost
		// 如果提供了任务ID但无法获取到模型，且是消费或退款类型，需要提供花费金额
		if req.TaskId != "" && (req.Type == usageDetailModel.UsageDetailTypeConsumption || req.Type == usageDetailModel.UsageDetailTypeRefund) && modelName == "" {
			if actualCost <= 0 {
				return errors.New("任务中没有模型信息，请提供花费金额")
			}
		} else if actualCost <= 0 {
			return errors.New("花费金额必须大于0")
		}
	}

	// 生成details详情文案
	details := req.Details
	// 如果是退款类型且提供了任务ID，自动生成详情文案
	if req.Type == usageDetailModel.UsageDetailTypeRefund && req.TaskId != "" {
		details = fmt.Sprintf("任务失败: %s, 退回额度: $%.6f", req.TaskId, actualCost)
	}

	usageDetail := &usageDetailModel.UsageDetail{
		UserId:    userId,
		Type:      req.Type, // 直接存储英文类型
		ModelName: modelName,
		Cost:      actualCost,
		Details:   details,
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

	// 将类型转换为英文类型（如果传入的是中文类型，转换为英文进行查询）
	typeForQuery := usageDetailModel.ConvertTypeToEN(req.Type)

	list, err := usageDetailRepo.GetUsageDetailList(userId, req.Page, req.PageSize, req.StartTime, req.EndTime, typeForQuery, req.Model)
	if err != nil {
		log.Printf("usageDetailService.GetUsageDetailList 获取使用明细列表失败, error: %v", err)
		return nil, 0, errors.New("获取使用明细列表失败")
	}

	count, err := usageDetailRepo.CountUsageDetailList(userId, req.StartTime, req.EndTime, typeForQuery, req.Model)
	if err != nil {
		log.Printf("usageDetailService.GetUsageDetailList 统计使用明细列表总数失败, error: %v", err)
		return nil, 0, errors.New("统计使用明细列表总数失败")
	}

	return list, count, nil
}
