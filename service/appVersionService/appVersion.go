package appVersionService

import (
	"jiyu/model/appVersionModel"
	"jiyu/repo/appVersionRepo"
	"strings"

	"github.com/spf13/cast"
)

func GetVersion(platform string, buildVersion string) (appVersionModel.AppVersion, error) {
	result, err := appVersionRepo.Find(strings.ToLower(platform))
	if err != nil {
		return result, err
	}

	result.Status = false
	if cast.ToInt(buildVersion) < cast.ToInt(result.BuildVersion) {
		result.Status = true
	}

	return result, nil
}
