package usageDetailRepo

import (
	"jiyu/global"
	"jiyu/model/usageDetailModel"
	"time"
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

