package feedbackRepo

import (
	"jiyu/global"
	"jiyu/model/globalModel"
)

func SaveFeedback(feedback globalModel.Feedback) error {
	return global.DB.Save(&feedback).Error
}

func GetFeedback(userId, feedbackId int) (globalModel.Feedback, error) {
	var feedback globalModel.Feedback
	err := global.DB.Model(&globalModel.Feedback{}).
		Where("user_id = ?", userId).
		Where("id = ?", feedbackId).
		First(&feedback).Error
	return feedback, err
}

func ReadFeedback(userId, feedbackId int) error {
	return global.DB.Model(&globalModel.Feedback{}).
		Where("user_id = ?", userId).
		Where("replay != ?", "").
		Where("id = ?", feedbackId).
		Update("is_read", true).Error
}

func IsNewFeedbackReply(userId int) (bool, error) {
	var count int64
	err := global.DB.Model(&globalModel.Feedback{}).
		Where("user_id = ?", userId).
		Where("replay != ?", "").
		Where("is_read = ?", false).
		Count(&count).Error
	return count > 0, err
}

func GetList(userId int, page, pageSize int) ([]globalModel.Feedback, error) {
	var results []globalModel.Feedback
	err := global.DB.Model(&globalModel.Feedback{}).
		Where("user_id = ?", userId).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&results).Error
	return results, err
}

func Count(userId int) (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.Feedback{}).
		Where("user_id = ?", userId).
		Count(&count).Error
	return int(count), err
}

func GetQaList(page, pageSize int) ([]globalModel.QAListResponse, error) {
	var results []globalModel.QAListResponse
	err := global.DB.Model(&globalModel.QA{}).
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&results).Error
	return results, err
}

func CountQa() (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.QA{}).Count(&count).Error
	return int(count), err
}
