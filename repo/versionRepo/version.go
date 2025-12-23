package versionrepo

import (
	"jiyu/global"
	versionmodel "jiyu/model/versionModel"
)

func Updates(version *versionmodel.AppVersion) error {

	iosData := map[string]interface{}{"build_version": version.Ios}
	andData := map[string]interface{}{"build_version": version.Android}

	err := global.DB.Debug().Model(&versionmodel.AppVersion{}).Where("platform = ?", "ios").Updates(iosData).Error
	if err != nil {
		return err
	}
	err = global.DB.Debug().Model(&versionmodel.AppVersion{}).Where("platform = ?", "android").Updates(andData).Error

	return err
}
