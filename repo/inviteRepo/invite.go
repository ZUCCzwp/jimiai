package inviteRepo

import (
	"jiyu/global"
	"jiyu/model/inviteModel"
)

func ListByUserID(userId int64, page, pageSize int) ([]inviteModel.InviteLog, error) {
	var results []inviteModel.InviteLog
	err := global.DB.Model(&inviteModel.InviteLog{}).
		Where("user_id = ?", userId).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&results).Error
	return results, err
}

func AdminList(page int, pageSize int) ([]inviteModel.InviteLog, error) {
	var results []inviteModel.InviteLog
	err := global.DB.Model(&inviteModel.InviteLog{}).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&results).Error
	return results, err
}

func AdminInviteSuperior(uid int, page int, limit int) ([]inviteModel.InviteLog, error) {
	var results []inviteModel.InviteLog
	err := global.DB.Model(&inviteModel.InviteLog{}).
		Where("user_id = ?", uid).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&results).Error
	return results, err
}

func AdminInviteSuperiorCount(uid int) (int64, error) {
	var count int64
	err := global.DB.Model(&inviteModel.InviteLog{}).
		Where("user_id = ?", uid).
		Count(&count).Error
	return count, err
}

func AdminInviteInferior(uid int, page int, limit int) ([]inviteModel.InviteLog, error) {
	var results []inviteModel.InviteLog
	err := global.DB.Model(&inviteModel.InviteLog{}).
		Where("invite_user_id = ?", uid).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&results).Error
	return results, err
}

func AdminInviteInferiorCount(uid int) (int64, error) {
	var count int64
	err := global.DB.Model(&inviteModel.InviteLog{}).
		Where("invite_user_id = ?", uid).
		Count(&count).Error
	return count, err
}

func AdminCount() (int64, error) {
	var count int64
	err := global.DB.Model(&inviteModel.InviteLog{}).Count(&count).Error
	return count, err
}
