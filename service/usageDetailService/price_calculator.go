package usageDetailService

// ModelPriceConfig 模型价格配置
type ModelPriceConfig struct {
	ModelName string  // 模型名称
	Price     float64 // 普通用户价格
}

// GetModelPrice 获取模型价格配置
// 这里可以根据实际需求从配置文件或数据库读取
var modelPriceMap = map[string]float64{
	"sora_url":                   0.01, // 普通用户价格
	"gpt-4":                      0.1,  // 普通用户价格
	"claude-3":                   0.09, // 普通用户价格
	"gpt-3.5":                    0.09, // 普通用户价格
	"dall-e":                     1.5,  // 普通用户价格
	"midjourney":                 1.5,  // 普通用户价格
	"sora2-portrait":             0.09, // 普通用户价格
	"sora2-landscape":            0.09, // 普通用户价格
	"sora2-portrait-15s":         0.09, // 普通用户价格
	"sora2-landscape-15s":        0.09, // 普通用户价格
	"sora2-pro-portrait-25s":     1.5,  // 普通用户价格
	"sora2-pro-landscape-25s":    1.5,  // 普通用户价格
	"sora2-pro-portrait-hd-15s":  1.5,  // 普通用户价格
	"sora2-pro-landscape-hd-15s": 1.5,  // 普通用户价格
}

// VipDiscountRate VIP用户折扣率（0.85 表示85%，即15%折扣）
const VipDiscountRate = 0.85

// CalculateCost 根据模型名称和用户VIP状态计算实际花费
func CalculateCost(modelName string, isVip bool) float64 {
	// 获取模型的基础价格
	basePrice, exists := modelPriceMap[modelName]
	if !exists {
		// 如果模型不存在，返回默认价格
		basePrice = 0.01
	}

	// 如果是VIP用户，应用折扣
	if isVip {
		return basePrice * VipDiscountRate
	}

	return basePrice
}

// GetModelPrice 获取模型价格（不应用VIP折扣）
func GetModelPrice(modelName string) float64 {
	price, exists := modelPriceMap[modelName]
	if !exists {
		return 0.01 // 默认价格
	}
	return price
}

// GetVipModelPrice 获取VIP用户的模型价格
func GetVipModelPrice(modelName string) float64 {
	return CalculateCost(modelName, true)
}
