package usageDetailRepo

import (
	"jiyu/global"
	"jiyu/model/usageDetailModel"
	"time"

	"github.com/shopspring/decimal"
)

// CreateUsageDetail 创建使用明细
func CreateUsageDetail(usageDetail *usageDetailModel.UsageDetail) error {
	return global.DB.Model(&usageDetailModel.UsageDetail{}).Create(usageDetail).Error
}

// GetUsageDetailList 获取使用明细列表
func GetUsageDetailList(userId int, page, pageSize int, startTime, endTime, detailType, model string) ([]usageDetailModel.UsageDetailListResponse, error) {
	var results []usageDetailModel.UsageDetailListResponse

	query := global.DB.Model(&usageDetailModel.UsageDetail{}).
		Select("id, created_at as time, type, model_name as model, cost, details").
		Where("user_id = ?", userId)

	// 时间范围筛选（最多7天）
	if startTime != "" && endTime != "" {
		start, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err == nil {
			end, err := time.Parse("2006-01-02 15:04:05", endTime)
			if err == nil {
				// 检查时间范围是否超过7天
				if end.Sub(start).Hours() <= 168 { // 7天 = 168小时
					query = query.Where("created_at >= ? AND created_at <= ?", start, end)
				}
			}
		}
	} else {
		// 如果没有指定时间范围，默认只查询最近7天的记录
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		query = query.Where("created_at >= ?", sevenDaysAgo)
	}

	// 类型筛选
	if detailType != "" && detailType != "全部" {
		query = query.Where("type = ?", detailType)
	}

	// 模型名称筛选
	if model != "" {
		query = query.Where("model_name LIKE ?", "%"+model+"%")
	}

	// 分页和排序
	err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&results).Error

	return results, err
}

// CountUsageDetailList 统计使用明细列表总数
func CountUsageDetailList(userId int, startTime, endTime, detailType, model string) (int, error) {
	var count int64

	query := global.DB.Model(&usageDetailModel.UsageDetail{}).
		Where("user_id = ?", userId)

	// 时间范围筛选（最多7天）
	if startTime != "" && endTime != "" {
		start, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err == nil {
			end, err := time.Parse("2006-01-02 15:04:05", endTime)
			if err == nil {
				// 检查时间范围是否超过7天
				if end.Sub(start).Hours() <= 168 { // 7天 = 168小时
					query = query.Where("created_at >= ? AND created_at <= ?", start, end)
				}
			}
		}
	} else {
		// 如果没有指定时间范围，默认只查询最近7天的记录
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		query = query.Where("created_at >= ?", sevenDaysAgo)
	}

	// 类型筛选
	if detailType != "" && detailType != "全部" {
		query = query.Where("type = ?", detailType)
	}

	// 模型名称筛选
	if model != "" {
		query = query.Where("model_name LIKE ?", "%"+model+"%")
	}

	err := query.Count(&count).Error
	return int(count), err
}

// GetUsedCoinTotal 计算用户已用额度（消费类型为正数累加，退款类型为负数累加）
// 已用额度 = 所有消费类型的cost之和 - 所有退款类型的cost之和
// 使用decimal进行精确计算，避免浮点数精度丢失
func GetUsedCoinTotal(userId int) (decimal.Decimal, error) {
	var result struct {
		Total float64
	}

	// 计算消费类型的总额（正数）
	err := global.DB.Model(&usageDetailModel.UsageDetail{}).
		Select("COALESCE(SUM(cost), 0) as total").
		Where("user_id = ? AND type = ?", userId, usageDetailModel.UsageDetailTypeConsumption).
		Scan(&result).Error

	if err != nil {
		return decimal.Zero, err
	}

	consumptionTotal := decimal.NewFromFloat(result.Total)

	// 计算退款类型的总额（正数）
	err = global.DB.Model(&usageDetailModel.UsageDetail{}).
		Select("COALESCE(SUM(cost), 0) as total").
		Where("user_id = ? AND type = ?", userId, usageDetailModel.UsageDetailTypeRefund).
		Scan(&result).Error

	if err != nil {
		return decimal.Zero, err
	}

	refundTotal := decimal.NewFromFloat(result.Total)

	// 已用额度 = 消费总额 - 退款总额（使用decimal进行精确计算）
	usedCoin := consumptionTotal.Sub(refundTotal)

	return usedCoin, nil
}

// GetConsumptionCount 统计用户的所有消费记录总数（请求次数）
func GetConsumptionCount(userId int) (int, error) {
	var count int64

	err := global.DB.Model(&usageDetailModel.UsageDetail{}).
		Where("user_id = ? AND type = ?", userId, usageDetailModel.UsageDetailTypeConsumption).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
