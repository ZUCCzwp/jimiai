package chargeRepo

import (
	"jiyu/global"
	"jiyu/model/chargeModel"
)

func UpdateChargeConfig(chargeConfig chargeModel.ChargeInfo) error {
	return global.DB.Model(&chargeModel.ChargeInfo{}).
		Where("month = ?", chargeConfig.Month).
		Updates(chargeConfig).Error
}

func ListChargeConfig() ([]chargeModel.ChargeInfo, error) {
	var chargeConfig []chargeModel.ChargeInfo

	err := global.DB.Model(&chargeModel.ChargeInfo{}).
		Order("month asc").
		Find(&chargeConfig).Error
	return chargeConfig, err
}

func ListPageChargeConfig(page, limit int) ([]chargeModel.ChargeInfo, error) {
	var chargeConfig []chargeModel.ChargeInfo
	err := global.DB.Model(&chargeModel.ChargeInfo{}).
		Order("month asc").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&chargeConfig).Error
	return chargeConfig, err
}

func CountChargeConfig() (int, error) {
	var count int64
	err := global.DB.Model(&chargeModel.ChargeInfo{}).Count(&count).Error
	return int(count), err
}

func GetChargeConfigByMonth(month int) (chargeModel.ChargeInfo, error) {
	var chargeConfig chargeModel.ChargeInfo
	err := global.DB.Model(&chargeModel.ChargeInfo{}).
		Where("month = ?", month).
		First(&chargeConfig).Error
	return chargeConfig, err
}
