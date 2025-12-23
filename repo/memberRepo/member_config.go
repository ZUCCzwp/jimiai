package memberRepo

import (
	"jiyu/global"
	"jiyu/model/memberModel"
)

func ListMemberConfig() ([]memberModel.MemberConfig, error) {
	memberList := []memberModel.MemberConfig{}
	err := global.DB.Model(&memberModel.MemberConfig{}).Find(&memberList).Error
	if err != nil {
		return nil, err
	}
	return memberList, nil
}

func ListPageMemberConfig(page, limit int) ([]memberModel.MemberConfig, error) {
	memberList := []memberModel.MemberConfig{}
	err := global.DB.Model(&memberModel.MemberConfig{}).
		Order("id desc").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&memberList).Error
	if err != nil {
		return nil, err
	}
	return memberList, nil
}

func CountMemberConfig() (int64, error) {
	var count int64
	err := global.DB.Model(&memberModel.MemberConfig{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func UpdateMemberConfig(member memberModel.MemberConfig) error {
	updateMap := make(map[string]interface{})
	updateMap["type"] = member.Type
	updateMap["icon"] = member.Icon
	updateMap["name"] = member.Name
	updateMap["min_days"] = member.MinDays
	updateMap["max_days"] = member.MaxDays

	return global.DB.Model(&memberModel.MemberConfig{}).
		Where("type = ?", member.Type).
		Updates(updateMap).Error
}
