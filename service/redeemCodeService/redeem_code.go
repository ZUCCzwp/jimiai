package redeemCodeService

import (
	"errors"
	"jiyu/config"
	"jiyu/global"
	"jiyu/model/contextModel"
	"jiyu/model/redeemCodeModel"
	"jiyu/repo/redeemCodeRepo"
	"jiyu/repo/userRepo"
	"jiyu/util"
	"log"

	"gorm.io/gorm"
)

// CreateRedeemCode 创建兑换码
func CreateRedeemCode(req redeemCodeModel.CreateRedeemCodeRequest) (*redeemCodeModel.RedeemCodeResponse, error) {
	// 生成唯一的兑换码
	maxRetries := 10 // 最大重试次数
	var (
		code   string
		amount = req.Amount
		err    error
	)

	for i := 0; i < maxRetries; i++ {
		// 生成12位随机兑换码（字母+数字）
		code = util.GenerateRandomString(12)

		// 检查兑换码是否已存在
		err = redeemCodeRepo.ExistByCode(code)
		if err != nil {
			// 如果err不为nil，说明兑换码不存在，可以使用
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			// 如果是其他错误，记录日志并继续重试
			log.Printf("redeemCodeService.CreateRedeemCode 检查兑换码失败, error: %v", err)
		}
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("生成兑换码失败，请重试")
	}

	// 创建兑换码
	redeemCode := &redeemCodeModel.RedeemCode{
		Code:   code,
		Amount: amount,
		Status: redeemCodeModel.RedeemCodeStatusUnused,
	}

	err = redeemCodeRepo.CreateRedeemCode(redeemCode)
	if err != nil {
		log.Printf("redeemCodeService.CreateRedeemCode 创建兑换码失败, error: %v", err)
		return nil, errors.New("创建兑换码失败")
	}

	return &redeemCodeModel.RedeemCodeResponse{
		Code:      redeemCode.Code,
		Amount:    redeemCode.Amount,
		Status:    redeemCode.Status,
		CreatedAt: redeemCode.CreatedAt,
	}, nil
}

// RedeemCode 兑换兑换码
func RedeemCode(ctx contextModel.Context, req redeemCodeModel.RedeemCodeRequest) error {
	// 查找兑换码
	redeemCode, err := redeemCodeRepo.FindByCode(req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("兑换码不存在")
		}
		log.Printf("redeemCodeService.RedeemCode 查找兑换码失败, error: %v", err)
		return errors.New("查找兑换码失败")
	}

	// 检查兑换码是否已使用
	if redeemCode.Status == redeemCodeModel.RedeemCodeStatusUsed {
		return errors.New("兑换码已被使用")
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新兑换码状态
	userId := int(ctx.User.ID)
	err = redeemCodeRepo.UpdateRedeemCode(tx, req.Code, userId)
	if err != nil {
		tx.Rollback()
		log.Printf("redeemCodeService.RedeemCode 更新兑换码失败, error: %v", err)
		return errors.New("兑换失败")
	}

	// 获取用户信息
	user, err := userRepo.FindByID(userId)
	if err != nil {
		tx.Rollback()
		log.Printf("redeemCodeService.RedeemCode 获取用户信息失败, error: %v", err)
		return errors.New("兑换失败")
	}

	// 更新用户充值额度（JimiCoin）
	user.JimiCoin = redeemCode.Amount
	err = userRepo.UpdateJimiCoinTx(tx, user)
	if err != nil {
		tx.Rollback()
		log.Printf("redeemCodeService.RedeemCode 更新用户充值额度失败, error: %v", err)
		return errors.New("兑换失败")
	}

	// 判断用户 api_key 是否为空，如果为空则更新为配置文件中的 api_key
	if user.ApiKey == "" && config.AppConfig.ApiKey != "" {
		user.ApiKey = config.AppConfig.ApiKey
		err = userRepo.UpdateApiKeyTx(tx, user)
		if err != nil {
			tx.Rollback()
			log.Printf("redeemCodeService.RedeemCode 更新用户api_key失败, error: %v", err)
			return errors.New("兑换失败")
		}
	}

	// 提交事务
	tx.Commit()

	return nil
}

// GetRedeemCodeList 获取充值明细列表
func GetRedeemCodeList(page, pageSize int, code string, status int) ([]redeemCodeModel.RedeemCodeListResponse, int, error) {
	list, err := redeemCodeRepo.GetList(page, pageSize, code, status)
	if err != nil {
		log.Printf("redeemCodeService.GetRedeemCodeList 获取充值明细列表失败, error: %v", err)
		return nil, 0, errors.New("获取充值明细列表失败")
	}

	count, err := redeemCodeRepo.CountList(code, status)
	if err != nil {
		log.Printf("redeemCodeService.GetRedeemCodeList 统计充值明细列表总数失败, error: %v", err)
		return nil, 0, errors.New("统计充值明细列表总数失败")
	}

	return list, count, nil
}
