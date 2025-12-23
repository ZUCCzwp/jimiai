package versionservice

import (
	"jiyu/model/contextModel"
	versionmodel "jiyu/model/versionModel"
	versionrepo "jiyu/repo/versionRepo"
)

func UpdateVersion(ctx contextModel.Context, data versionmodel.AppVersion) error {

	err := versionrepo.Updates(&data)

	return err
}
