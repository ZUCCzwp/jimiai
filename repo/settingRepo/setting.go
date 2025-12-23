package settingRepo

import (
	"jiyu/global"
	"jiyu/model/settingModel"
)

func Find() (*settingModel.Setting, error) {
	var setting settingModel.Setting
	err := global.DB.Model(&settingModel.Setting{}).First(&setting).Error
	return &setting, err
}

func Updates(setting *settingModel.Setting) error {
	return global.DB.Model(&settingModel.Setting{}).
		Where("id = 1").
		Updates(&setting).Error
}

func Save(setting *settingModel.Setting) error {
	return global.DB.Save(setting).Error
}

func Count() (int, error) {
	var count int64
	err := global.DB.Model(&settingModel.Setting{}).Where("1 = 1").Count(&count).Error
	return int(count), err
}
