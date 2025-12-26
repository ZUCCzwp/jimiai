package userRepo

import (
	"context"
	"fmt"
	"jiyu/global"
	"jiyu/model/userModel"
	"strings"
	"time"

	"gorm.io/gorm"
)

func Create(user *userModel.User) (uint, error) {
	err := global.DB.Model(&userModel.User{}).Create(user).Error

	return user.ID, err
}

func Find(phone string) (*userModel.User, error) {
	var user userModel.User

	err := global.DB.Model(&userModel.User{}).Where("phone = ?", phone).First(&user).Error

	return &user, err
}

func FindByID(id int) (*userModel.User, error) {
	var user userModel.User

	err := global.DB.Model(&userModel.User{}).Where("id = ?", id).First(&user).Error

	return &user, err
}

// FindByEmail 根据邮箱查找用户
func FindByEmail(email string) (*userModel.User, error) {
	var user userModel.User
	email = strings.ToLower(strings.TrimSpace(email))

	err := global.DB.Model(&userModel.User{}).Where("email = ?", email).First(&user).Error

	return &user, err
}

// FindByMyInviteCode 根据自身邀请码查找用户
func FindByMyInviteCode(inviteCode string) (*userModel.User, error) {
	var user userModel.User
	inviteCode = strings.TrimSpace(inviteCode)

	err := global.DB.Model(&userModel.User{}).Where("my_invite_code = ?", inviteCode).First(&user).Error

	return &user, err
}

func FindUserByIds(ids []string) ([]userModel.User, error) {
	var results []userModel.User
	err := global.DB.Model(&userModel.User{}).Where("id In ?", ids).Find(&results).Error
	return results, err
}

func FindUserByIntervalDay(day int) (int, error) {
	var count int64

	err := global.DB.Model(&userModel.User{}).Where("created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", day).Count(&count).Error

	return int(count), err
}

// FindUserByLsatRequestTime 根据最近上线时间先后获取用户
func FindUserByLsatRequestTime(count int) ([]userModel.User, error) {
	var users []userModel.User

	err := global.DB.Model(&userModel.User{}).Order("last_request_time DESC").Limit(count).Find(&users).Error

	return users, err
}

// FindByUserName 根据用户名获取用户
func FindByUserName(username string) (*userModel.User, error) {
	var user userModel.User

	err := global.DB.Model(&userModel.User{}).Where("nickname = ?", username).First(&user).Error

	return &user, err
}

func ListAll() ([]userModel.User, error) {
	var users []userModel.User

	err := global.DB.Model(&userModel.User{}).Find(&users).Error

	return users, err
}

// IncreaseRefCountById 累加用户的邀请次数
func IncreaseRefCountById(uid int) error {
	return global.DB.Model(&userModel.User{}).Where("id = ?", uid).Update("ref_count", gorm.Expr("ref_count + ?", 1)).Error
}

func Delete(id int) error {
	// 删除用户
	fmt.Println(id)
	// return nil
	err := global.DB.Unscoped().Where("id = ?", id).Delete(&userModel.User{}).Error
	return err
}

func List(page, limit int, uid int, phone, nickname string) ([]userModel.User, error) {
	var users []userModel.User

	where := make([]string, 0)
	if uid != 0 {
		where = append(where, fmt.Sprintf("id = %d", uid))
	}

	if phone != "" {
		where = append(where, fmt.Sprintf("phone LIKE '%s'", "%"+phone+"%"))
	}

	if nickname != "" {
		where = append(where, fmt.Sprintf("info_nickname LIKE '%s'", "%"+nickname+"%"))
	}

	wheres := strings.Join(where, " AND ")

	err := global.DB.Model(&userModel.User{}).
		Where(wheres).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&users).Error
	return users, err
}

func Count(phone, nickname string) (int, error) {
	var count int64
	where := make([]string, 0)
	if phone != "" {
		where = append(where, fmt.Sprintf("phone LIKE '%s'", "%"+phone+"%"))
	}

	if nickname != "" {
		where = append(where, fmt.Sprintf("info_nickname LIKE '%s'", "%"+nickname+"%"))
	}

	wheres := strings.Join(where, " AND ")
	err := global.DB.Model(&userModel.User{}).Where(wheres).Count(&count).Error
	return int(count), err
}

func Exist(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	return global.DB.Model(&userModel.User{}).Where("email = ?", email).First(&userModel.User{}).Error
}

// ExistByNickname 检查昵称（账号）是否存在
func ExistByNickname(nickname string) error {
	nickname = strings.TrimSpace(nickname)
	return global.DB.Model(&userModel.User{}).Where("nickname = ?", nickname).First(&userModel.User{}).Error
}

// ExistByMyInviteCode 检查自身邀请码是否存在
func ExistByMyInviteCode(inviteCode string) error {
	inviteCode = strings.TrimSpace(inviteCode)
	return global.DB.Model(&userModel.User{}).Where("my_invite_code = ?", inviteCode).First(&userModel.User{}).Error
}

func Save(user *userModel.User) error {
	return global.DB.Save(user).Error
}

func SaveWithCtx(ctx context.Context, user *userModel.User) error {
	return global.GetDB(ctx).Save(user).Error
}

func Updates(user *userModel.User) error {
	return global.DB.Debug().
		Model(&userModel.User{}).
		Where("id = ?", user.ID).
		Updates(user).Error
}

func UpdateUserInfo(user *userModel.User) error {
	updateData := map[string]interface{}{
		"city": user.City,
	}
	return global.DB.Debug().
		Model(&userModel.User{}).
		Where("id = ?", user.ID).
		Updates(updateData).Error
}

func UpdateUserLike(user *userModel.User) error {
	return global.DB.Debug().
		Model(&userModel.User{}).
		Select("like_shape").
		Where("id = ?", user.ID).
		Updates(user).Error
}

func UpdatePassword(user *userModel.User) error {
	return global.DB.Debug().
		Model(&userModel.User{}).
		Select("password").
		Where("id = ?", user.ID).
		Updates(user).Error
}

func UpdateApiKeyTx(tx *gorm.DB, user *userModel.User) error {
	return tx.Model(&userModel.User{}).
		Where("id = ?", user.ID).
		Update("api_key", user.ApiKey).Error
}

func UpdateInviteCountTx(tx *gorm.DB, user *userModel.User) error {
	return tx.Model(&userModel.User{}).
		Where("id = ?", user.ID).
		Update("ref_count", gorm.Expr("ref_count + ?", 1)).Error
}

func UpdateJimiCoinTx(tx *gorm.DB, user *userModel.User) error {
	return tx.Model(&userModel.User{}).
		Where("id = ?", user.ID).
		Update("jimi_coin", gorm.Expr("jimi_coin + ?", user.JimiCoin)).Error
}

func UpdateVipExpireTimeTx(tx *gorm.DB, user *userModel.User) error {
	return tx.Model(&userModel.User{}).
		Select("vip_time").
		Where("id = ?", user.ID).
		Update("vip_time", user.VipTime).Error
}

func SaveTx(tx *gorm.DB, user *userModel.User) error {
	return tx.Save(user).Error
}

func OnlineUserCount() int {
	var count int64

	global.DB.Model(&userModel.User{}).
		Where("last_request_time BETWEEN DATE_SUB(NOW(), INTERVAL 15 MINUTE)").
		Count(&count)

	return int(count)
}

func UpdateUserHotValue(tx *gorm.DB, uid, value int) error {
	err := tx.Model(&userModel.User{}).Debug().Where("id = ?", uid).Update("hot_value_count", gorm.Expr("hot_value_count + ?", value)).Error
	return err
}

// SelectHotUser 选择当前符合条件的热门用户
func SelectHotUser(shapeType, count int) ([]userModel.HotUser, error) {
	var users []userModel.HotUser

	err := global.DB.
		Model(&userModel.User{}).
		Select("id as uid, info_shape as shape, phone, info_nickname as nickname, info_avatar as avatar_url, hot_value_count, hot_const_value").
		Where("info_shape = ? and hot_value_count > 0", shapeType).
		Order("hot_value_count+hot_const_value DESC").
		Limit(count).
		Find(&users).
		Error

	return users, err
}

func SaveHotUsers(data []userModel.HotUser) error {
	return global.DB.Save(data).Error
}

func FindHotUsers() ([]userModel.HotUser, error) {
	var users []userModel.HotUser
	err := global.DB.Model(&userModel.HotUser{}).Where("1=1").Find(&users).Error
	return users, err
}

func DropHotUsers() error {
	return global.DB.Unscoped().Where("1=1").Delete(&userModel.HotUser{}).Error
}

func DeleteByID(id int) error {
	return global.DB.Unscoped().Model(&userModel.User{}).Where("id = ?", id).Delete(&userModel.User{}).Error
}

func CountMember() (int64, error) {
	var count int64
	err := global.DB.Model(&userModel.User{}).
		Where("vip_time > ?", time.Now()).
		Count(&count).Error
	return count, err
}

func ListMemberByPage(page, limit int) ([]userModel.User, error) {
	var users []userModel.User
	err := global.DB.Model(&userModel.User{}).
		Where("vip_time > ?", time.Now()).
		Offset((page - 1) * limit).
		Limit(limit).Find(&users).Error
	return users, err
}

func EditMemberVipTime(uid int, member *userModel.User) error {
	updateMap := make(map[string]any)
	updateMap["vip_time"] = member.VipTime
	return global.DB.Model(&userModel.User{}).
		Where("id = ?", uid).
		Updates(updateMap).Error
}

func UpdateCount(uid, push, slvage int) error {
	var err error
	if push > 0 {
		err = global.DB.Model(&userModel.User{}).Where("id = ?", uid).Update("push_bottle_count", gorm.Expr("push_bottle_count + ?", push)).Error
	} else if slvage > 0 {
		err = global.DB.Model(&userModel.User{}).Where("id = ?", uid).Update("salvage_count", gorm.Expr("salvage_count + ?", slvage)).Error
	}
	return err
}
