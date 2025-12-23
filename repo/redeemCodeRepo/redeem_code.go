package redeemCodeRepo

import (
	"time"

	"jiyu/config"
	"jiyu/global"
	"jiyu/model/redeemCodeModel"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CreateRedeemCode 创建兑换码
func CreateRedeemCode(redeemCode *redeemCodeModel.RedeemCode) error {
	return global.DB.Model(&redeemCodeModel.RedeemCode{}).Create(redeemCode).Error
}

// FindByCode 根据兑换码查找
func FindByCode(code string) (*redeemCodeModel.RedeemCode, error) {
	var redeemCode redeemCodeModel.RedeemCode
	err := global.DB.Model(&redeemCodeModel.RedeemCode{}).
		Where("code = ?", code).
		First(&redeemCode).Error
	return &redeemCode, err
}

// ExistByCode 检查兑换码是否存在
func ExistByCode(code string) error {
	return global.DB.Model(&redeemCodeModel.RedeemCode{}).
		Where("code = ?", code).
		First(&redeemCodeModel.RedeemCode{}).Error
}

// UpdateRedeemCode 更新兑换码（兑换时使用）
func UpdateRedeemCode(tx *gorm.DB, code string, userId int) error {
	now := time.Now()
	return tx.Model(&redeemCodeModel.RedeemCode{}).
		Where("code = ?", code).
		Updates(map[string]interface{}{
			"status":      redeemCodeModel.RedeemCodeStatusUsed,
			"user_id":     userId,
			"redeemed_at": &now,
		}).Error
}

// GetList 获取充值明细列表（关联用户表）
func GetList(page, pageSize int, code string, status int) ([]redeemCodeModel.RedeemCodeListResponse, error) {
	var results []redeemCodeModel.RedeemCodeListResponse

	// 使用GORM的NamingStrategy获取表名
	namingStrategy := schema.NamingStrategy{TablePrefix: config.DBConfig.TablePrefix}
	redeemCodeTable := namingStrategy.TableName("RedeemCode")
	userTable := namingStrategy.TableName("User")

	query := global.DB.Model(&redeemCodeModel.RedeemCode{}).
		Select(redeemCodeTable + ".code, " + redeemCodeTable + ".status, " + userTable + ".nickname as user_nickname, " + redeemCodeTable + ".redeemed_at, " + redeemCodeTable + ".amount").
		Joins("LEFT JOIN " + userTable + " ON " + redeemCodeTable + ".user_id = " + userTable + ".id")

	// 条件筛选
	if code != "" {
		query = query.Where(redeemCodeTable+".code LIKE ?", "%"+code+"%")
	}
	if status >= 0 {
		query = query.Where(redeemCodeTable+".status = ?", status)
	}

	// 分页和排序
	err := query.
		Order(redeemCodeTable + ".id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&results).Error

	return results, err
}

// CountList 统计充值明细列表总数
func CountList(code string, status int) (int, error) {
	var count int64

	query := global.DB.Model(&redeemCodeModel.RedeemCode{})

	// 条件筛选
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	return int(count), err
}
