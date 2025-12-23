package settingService

import (
	"jiyu/model/contextModel"
	"jiyu/model/settingModel"
	"jiyu/repo/settingRepo"
)

func Get(ctx contextModel.AdminContext) (*settingModel.Setting, error) {
	setting, err := settingRepo.Find()
	if err != nil {
		return &settingModel.Setting{}, err
	}

	return setting, nil
}

func Update(ctx contextModel.AdminContext, data settingModel.Setting) error {
	return settingRepo.Updates(&data)
}
