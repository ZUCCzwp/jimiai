package adminUserRepo

import (
	"jiyu/global"
	"jiyu/model/adminUserModel"
)

func List() ([]adminUserModel.User, error) {
	var users []adminUserModel.User
	err := global.DB.Model(&adminUserModel.User{}).Find(&users).Error
	return users, err
}

func Count() (int, error) {
	var count int64
	err := global.DB.Model(&adminUserModel.User{}).Count(&count).Error
	return int(count), err
}

func Create(user *adminUserModel.User) error {
	err := global.DB.Model(&adminUserModel.User{}).Create(user).Error
	return err
}

func Save(user *adminUserModel.User) error {
	err := global.DB.Save(user).Error
	return err
}

func Find(id int) (*adminUserModel.User, error) {
	var user adminUserModel.User
	err := global.DB.Model(&adminUserModel.User{}).Where("id = ?", id).First(&user).Error
	return &user, err
}

func Delete(id int) error {
	err := global.DB.Unscoped().Model(&adminUserModel.User{}).Where("id = ?", id).Delete(&adminUserModel.User{}).Error
	return err
}

func FindByUsername(username string) (*adminUserModel.User, error) {
	var user adminUserModel.User
	err := global.DB.Model(&adminUserModel.User{}).Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindUserAttrs() ([]adminUserModel.UserAttr, error) {
	var results []adminUserModel.UserAttr
	err := global.DB.Model(&adminUserModel.UserAttr{}).Order("sort_id").Find(&results).Error
	return results, err
}

func FindAttrByTypeAndKey(attrType, attrKey string) (adminUserModel.UserAttr, error) {
	var data adminUserModel.UserAttr
	err := global.DB.Model(&adminUserModel.UserAttr{}).Where("attr_type = ? and attr_key = ?", attrType, attrKey).First(&data).Error
	return data, err
}

func SaveUserAttrs(attrs []adminUserModel.UserAttr) error {
	err := global.DB.Model(&adminUserModel.UserAttr{}).Save(&attrs).Error
	return err
}

func CountUserAttr() (int, error) {
	var count int64
	err := global.DB.Model(&adminUserModel.UserAttr{}).Count(&count).Error
	return int(count), err
}

func CountUserTag() (int, error) {
	var count int64
	err := global.DB.Model(&adminUserModel.UserTag{}).Count(&count).Error
	return int(count), err
}

func FindUserTags(page, limit int) ([]adminUserModel.UserTag, error) {
	var results []adminUserModel.UserTag
	err := global.DB.Model(&adminUserModel.UserTag{}).Order("sort_id").Offset((page - 1) * limit).Limit(limit).Find(&results).Error
	return results, err
}

func SaveUserTag(tag *adminUserModel.UserTag) error {
	err := global.DB.Save(tag).Error
	return err
}

func DeleteUserTag(id int) error {
	// 删除tag
	err := global.DB.Unscoped().Model(&adminUserModel.UserTag{}).Where("id = ?", id).Delete(&adminUserModel.UserTag{}).Error
	return err
}
