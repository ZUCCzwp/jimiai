package adminReportRepo

import (
	"jiyu/global"
	"jiyu/model/globalModel"
)

func ListForAdmin(page, limit int) ([]globalModel.Report, error) {
	var results []globalModel.Report
	err := global.DB.Model(&globalModel.Report{}).Order("id desc").Offset((page - 1) * limit).Limit(limit).Find(&results).Error
	return results, err
}

func Count() (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.Report{}).Count(&count).Error
	return int(count), err
}

func UpdateBanLevel(id, level int) error {
	return global.DB.Model(&globalModel.Report{}).Where("id = ?", id).Update("ban_level", level).Error
}

func FindById(id int) (*globalModel.Report, error) {
	var data globalModel.Report
	err := global.DB.Model(&globalModel.Report{}).Where("id = ?", id).First(&data).Error
	return &data, err
}

func Save(r *globalModel.Report) error {
	return global.DB.Save(r).Error
}

func SaveQA(qa globalModel.QA) error {
	return global.DB.Save(&qa).Error
}

func QAListForAdmin(page, limit int) ([]globalModel.QA, error) {
	var results []globalModel.QA
	err := global.DB.Model(&globalModel.QA{}).Order("id desc").Offset((page - 1) * limit).Limit(limit).Find(&results).Error
	return results, err
}

func CountQA() (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.QA{}).Count(&count).Error
	return int(count), err
}

func FeedbackListForAdmin(page, limit, uid int) ([]globalModel.Feedback, error) {
	var results []globalModel.Feedback
	tb := global.DB.Model(&globalModel.Feedback{})
	if uid > 0 {
		tb = tb.Where("user_id = ?", uid)
	}

	err := tb.Order("id desc").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&results).Error
	return results, err
}

func SaveFeedback(feedback globalModel.Feedback) error {
	return global.DB.Save(&feedback).Error
}

func CountFeedback(uid int) (int, error) {
	var count int64
	tb := global.DB.Model(&globalModel.Feedback{})
	if uid > 0 {
		tb = tb.Where("user_id = ?", uid)
	}
	err := tb.Count(&count).Error
	return int(count), err
}
