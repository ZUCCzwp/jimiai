package main

import (
	"jiyu/config"
	"jiyu/global"
	"jiyu/model/usageDetailModel"
	"jiyu/repo/usageDetailRepo"
	"jiyu/repo/userRepo"
	"log"
	"time"
)

// 创建使用明细测试数据
// 使用方法: go run cmd/seed/main.go
func main() {
	// 初始化配置
	config.InitConfig()

	// 连接数据库
	_ = global.NewDB()

	// 自动迁移
	global.AutoMigrate()

	// 获取第一个用户（如果没有用户，需要先创建用户）
	users, err := userRepo.ListAll()
	if err != nil || len(users) == 0 {
		log.Fatalf("获取用户失败或没有用户，请先创建用户: %v", err)
	}

	user := users[0]
	userId := int(user.ID)
	isVip := user.IsVip()
	log.Printf("使用用户ID: %d (昵称: %s, VIP: %v)\n", userId, user.Nickname, isVip)

	// 价格计算函数（与service层保持一致）
	calculateCost := func(modelName string, isVip bool) float64 {
		modelPriceMap := map[string]float64{
			"sora_url":   0.01,
			"gpt-4":      0.1,
			"claude-3":   0.09,
			"gpt-3.5":    0.09,
			"dall-e":     1.5,
			"midjourney": 1.5,
		}
		vipDiscountRate := 0.85

		basePrice, exists := modelPriceMap[modelName]
		if !exists {
			basePrice = 0.01
		}

		if isVip {
			return basePrice * vipDiscountRate
		}
		return basePrice
	}

	// 创建测试数据（使用英文类型）
	testData := []usageDetailModel.UsageDetail{
		{
			UserId:    userId,
			Type:      usageDetailModel.UsageDetailTypeCharge,
			ModelName: "",
			Cost:      100.0000,
			Details:   "账户充值",
		},
		{
			UserId:    userId,
			Type:      usageDetailModel.UsageDetailTypeRefund,
			ModelName: "",
			Cost:      50.0000,
			Details:   "订单退款",
		},
		{
			UserId:    userId,
			Type:      usageDetailModel.UsageDetailTypeConsumption,
			ModelName: "sora_url",
			Cost:      calculateCost("sora_url", isVip),
			Details:   "-",
		},
		{
			UserId:    userId,
			Type:      usageDetailModel.UsageDetailTypeConsumption,
			ModelName: "gpt-4",
			Cost:      calculateCost("gpt-4", isVip),
			Details:   "-",
		},
		{
			UserId:    userId,
			Type:      usageDetailModel.UsageDetailTypeConsumption,
			ModelName: "claude-3",
			Cost:      calculateCost("claude-3", isVip),
			Details:   "-",
		},
	}

	// 插入测试数据
	for i, data := range testData {
		err := usageDetailRepo.CreateUsageDetail(&data)
		if err != nil {
			log.Printf("创建测试数据失败 [%d]: %v\n", i+1, err)
			continue
		}
		
		log.Printf("✓ 创建测试数据成功 [%d]: 类型=%s, 模型=%s, 花费=%.4f\n", i+1, data.Type, data.ModelName, data.Cost)
		time.Sleep(200 * time.Millisecond) // 短暂延迟，确保时间戳不同
	}

	log.Println("\n所有测试数据创建完成！")
}

