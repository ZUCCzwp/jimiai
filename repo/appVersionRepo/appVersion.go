package appVersionRepo

import (
	"jiyu/global"
	"jiyu/model/appVersionModel"
)

func Find(platform string) (appVersionModel.AppVersion, error) {
	var version appVersionModel.AppVersion
	err := global.DB.Where("platform = ?", platform).First(&version).Error
	if err != nil {
		return version, err
	}
	return version, nil
}
