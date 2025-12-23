package serverRepo

import (
	"jiyu/global"
	"jiyu/model/serverModel"
)

func FindConfig() *serverModel.Config {
	var config serverModel.Config
	global.DB.Model(&serverModel.Config{}).Where("id = 1").First(&config)
	return &config
}

func CreateConfig(config *serverModel.Config) {
	global.DB.Create(config)
}

func SaveConfig(c *serverModel.Config) error {
	return global.DB.Save(c).Error
}
