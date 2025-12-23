package adminRouterRepo

import (
	"jiyu/global"
	"jiyu/model/adminRouterModel"
)

// ExistById 根据id判断是否存在
func ExistById(id int) bool {
	var user adminRouterModel.Router

	global.DB.Model(&adminRouterModel.Router{}).Where("id = ?", id).First(&user)

	return user.ID > 0
}

func Save(router adminRouterModel.Router) error {
	return global.DB.Save(&router).Error
}

func Delete(id int) error {
	return global.DB.Unscoped().Delete(&adminRouterModel.Router{}, id).Error
}

func List() ([]adminRouterModel.Router, error) {
	var routers []adminRouterModel.Router
	err := global.DB.Model(&adminRouterModel.Router{}).Find(&routers).Error
	return routers, err
}
