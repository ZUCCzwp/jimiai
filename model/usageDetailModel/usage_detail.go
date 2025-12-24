package usageDetailModel

import (
	"time"

	"gorm.io/gorm"
)

// UsageDetailType 使用明细类型（API请求使用）
const (
	UsageDetailTypeConsumption = "consumption" // 消费
	UsageDetailTypeRefund      = "refund"      // 退款
	UsageDetailTypeCharge      = "charge"      // 充值
)

// UsageDetailTypeCN 使用明细类型（数据库存储使用中文）
const (
	UsageDetailTypeConsumptionCN = "消费" // 消费
	UsageDetailTypeRefundCN      = "退款" // 退款
	UsageDetailTypeChargeCN      = "充值" // 充值
)

// ConvertTypeToCN 将英文类型转换为中文类型
func ConvertTypeToCN(typeEn string) string {
	typeMap := map[string]string{
		UsageDetailTypeConsumption: UsageDetailTypeConsumptionCN,
		UsageDetailTypeRefund:      UsageDetailTypeRefundCN,
		UsageDetailTypeCharge:      UsageDetailTypeChargeCN,
	}
	if cnType, ok := typeMap[typeEn]; ok {
		return cnType
	}
	return typeEn // 如果不在映射中，返回原值
}

// UsageDetail 使用明细表
type UsageDetail struct {
	gorm.Model
	UserId    int     `json:"user_id" gorm:"type:int;not null;default:0;comment:'用户ID';index:idx_user_id"`        // 用户ID
	Type      string  `json:"type" gorm:"type:varchar(20);not null;default:'';comment:'类型:消费/充值';index:idx_type"` // 类型
	ModelName string  `json:"model" gorm:"type:varchar(100);default:'';comment:'模型名称';index:idx_model"`           // 模型名称
	Cost      float64 `json:"cost" gorm:"type:decimal(10,4);not null;default:0;comment:'花费金额'"`                   // 花费金额
	Details   string  `json:"details" gorm:"type:varchar(500);default:'';comment:'详情'"`                           // 详情
}

// CreateUsageDetailRequest 创建使用明细请求
type CreateUsageDetailRequest struct {
	Type    string  `json:"type" binding:"required"` // 类型：consumption(消费)/refund(退款)/charge(充值)
	Model   string  `json:"model"`                   // 模型名称
	Cost    float64 `json:"cost"`                    // 花费金额
	Details string  `json:"details"`                 // 详情
}

// UsageDetailListRequest 使用明细列表请求
type UsageDetailListRequest struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
	Type      string `json:"type"`       // 类型筛选
	Model     string `json:"model"`      // 模型名称筛选
}

// UsageDetailListResponse 使用明细列表响应
type UsageDetailListResponse struct {
	ID      uint      `json:"id"`      // 序号
	Time    time.Time `json:"time"`    // 时间
	Type    string    `json:"type"`    // 类型
	Model   string    `json:"model"`   // 模型
	Cost    float64   `json:"cost"`    // 花费
	Details string    `json:"details"` // 详情
}
