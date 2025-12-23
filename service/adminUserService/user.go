package adminUserService

import (
	"errors"
	"jiyu/model/adminUserModel"
	"jiyu/model/contextModel"
	"jiyu/model/userModel"
	"jiyu/repo/adminUserRepo"
	"jiyu/repo/userRepo"
	"jiyu/util"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(data adminUserModel.LoginRequest) (string, error) {
	user, err := adminUserRepo.FindByUsername(data.Username)
	if err != nil {
		log.Printf("登陆失败: %s", err)
		return "", errors.New("用户不存在")
	}
	if !user.ComparePwd(data.Password) {
		return "", errors.New("密码错误")
	}

	jwt := util.CreateJWT(user.ID, user.Username)
	return jwt, nil
}

func ChangePwd(ctx contextModel.AdminContext, data adminUserModel.ChangePwdRequest) error {
	user, err := adminUserRepo.Find(int(ctx.User.ID))
	if err != nil {
		return errors.New("用户不存在")
	}
	if !user.ComparePwd(data.OldPwd) {
		log.Printf("密码错误: %s, pwd: %s", err, data.NewPwd)
		return errors.New("密码错误")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.NewPwd), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}
	user.Password = string(hash)
	err = adminUserRepo.Save(user)
	if err != nil {
		log.Printf("密码修改失败: %s", err)
		return errors.New("密码修改失败")
	}
	return nil
}

func List() ([]adminUserModel.UserResponse, error) {
	users, err := adminUserRepo.List()
	if err != nil {
		return nil, err
	}

	var results []adminUserModel.UserResponse
	for _, user := range users {
		role := "管理员"
		if user.Role == 1 {
			role = "超级管理员"
		}
		results = append(results, adminUserModel.UserResponse{
			Id:       int(user.ID),
			Username: user.Username,
			Role:     role,
		})
	}

	return results, nil
}

func AddUser(data adminUserModel.AddUserRequest) error {
	user := adminUserModel.NewUser()

	user.Username = data.Username
	user.Password = data.Password
	user.Role = data.Role
	user.HashPassword()

	err := adminUserRepo.Save(user)
	if err != nil {
		log.Printf("adminUserService.AddUser 添加用户失败: %s", err)
		return err
	}

	return nil
}

func ResetPwd(data adminUserModel.ResetPwdRequest) error {
	user, err := adminUserRepo.Find(data.Id)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Password = data.Password
	user.HashPassword()

	err = adminUserRepo.Save(user)
	if err != nil {
		log.Printf("adminUserService.ResetPwd 密码修改失败: %s", err)
		return errors.New("密码修改失败")
	}

	return nil
}

func ResetAuthority(data adminUserModel.ResetAuthorityRequest) error {
	user, err := adminUserRepo.Find(data.Id)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Role = data.Role

	err = adminUserRepo.Save(user)
	if err != nil {
		log.Printf("adminUserService.ResetPwd 权限修改失败: %s", err)
		return errors.New("权限修改失败")
	}

	return nil
}

func DeleteAdmin(id int) error {
	err := adminUserRepo.Delete(id)
	if err != nil {
		log.Printf("adminUserService.Delete 用户删除失败: %s", err)
		return errors.New("用户删除失败")
	}
	return nil
}

func UserList(ctx contextModel.AdminContext, page, limit int, uid int, phone, nickname string) ([]adminUserModel.UserListResponse, int, error) {
	users, err := userRepo.List(page, limit, uid, phone, nickname)
	if err != nil {
		return make([]adminUserModel.UserListResponse, 0), 0, err
	}

	count, err := userRepo.Count(phone, nickname)
	if err != nil {
		return make([]adminUserModel.UserListResponse, 0), 0, err
	}

	var results []adminUserModel.UserListResponse

	for _, v := range users {
		var status int
		if v.BanLevel == 6 {
			status = 6
		} else if v.BanLevel > 0 && v.BanTime.After(time.Now()) {
			status = 1
		} else {
			status = 0
		}

		results = append(results, adminUserModel.UserListResponse{
			Id:             v.ID,
			Phone:          v.Phone,
			RegisterDevice: v.RegisterDevice,
			RegisterTime:   v.CreatedAt,
			LastLoginTime:  v.LastRequestTime,
			LastLoginIp:    v.LastLoginIp,
			Status:         status,
		})
	}

	return results, count, nil
}

func UserByUid(uid int) ([]adminUserModel.UserListResponse, int, error) {
	user, err := userRepo.FindByID(uid)
	if err != nil {
		return make([]adminUserModel.UserListResponse, 0), 0, err
	}

	count := 1

	users := make([]*userModel.User, 0)
	users = append(users, user)

	var results []adminUserModel.UserListResponse

	for _, v := range users {
		var status int

		switch status {
		case 0:
			status = 0
			// break
		case 1, 2, 3, 4, 5:
			status = 1
			// break
		case 6:
			status = 2
			// break
		}

		results = append(results, adminUserModel.UserListResponse{
			Id:             v.ID,
			Phone:          v.Phone,
			RegisterDevice: v.RegisterDevice,
			RegisterTime:   v.CreatedAt,
			LastLoginTime:  v.LastRequestTime,
			LastLoginIp:    v.LastLoginIp,
			Status:         v.BanLevel,
		})
	}

	return results, count, nil
}

func Ban(ctx contextModel.AdminContext, data adminUserModel.UserBanRequest) error {
	user, err := userRepo.FindByID(data.Uid)

	if err != nil {
		return errors.New("用户不存在")
	}

	user.Ban(data.BanLevel)

	err = userRepo.Save(user)
	if err != nil {
		return errors.New("用户禁言失败")
	}

	return nil
}

func Delete(ctx contextModel.AdminContext, uid int) error {
	err := userRepo.Delete(uid)
	if err != nil {
		log.Println("adminUserService.Delete 用户删除失败 err: ", err)
		return errors.New("用户删除失败")
	}
	return nil
}

func Detail(ctx contextModel.AdminContext, uid int) (adminUserModel.UserDetailResponse, error) {
	user, err := userRepo.FindByID(uid)
	if err != nil {
		log.Printf("adminUserService.Detail 用户查询失败 id: %d, err: %s", uid, err)
		return adminUserModel.UserDetailResponse{}, errors.New("用户查询失败")
	}

	result := adminUserModel.UserDetailResponse{
		Id:    int(user.ID),
		Phone: user.Phone,
	}
	return result, nil
}

func Update(ctx contextModel.AdminContext, data adminUserModel.UpdateUserRequest) error {

	user := userModel.User{
		Model: gorm.Model{ID: uint(data.Id)},
		Phone: data.Phone,
	}
	err := userRepo.Updates(&user)
	if err != nil {
		log.Printf("adminUserService.Update 用户更新失败 id: %d, err: %s", data.Id, err)
		return errors.New("用户更新失败")
	}
	return nil
}

func GetAttrs(ctx contextModel.AdminContext) ([]adminUserModel.UserAttrResponse, error) {
	results := make([]adminUserModel.UserAttrResponse, 0)

	attrs, err := adminUserRepo.FindUserAttrs()
	if err != nil {
		return results, err
	}

	for _, attr := range attrs {
		results = append(results, adminUserModel.UserAttrResponse{
			Id:        int(attr.ID),
			SortId:    attr.SortId,
			AttrType:  attr.AttrType,
			AttrName:  attr.AttrName,
			Color:     attr.Color,
			HaveRange: attr.HaveRange,
			Range:     attr.Range,
		})
	}

	return results, nil
}

func GetAttrByTypeAndKey(attrType, attrKey string) (adminUserModel.UserAttr, error) {
	return adminUserRepo.FindAttrByTypeAndKey(attrType, attrKey)
}

func UpdateAttr(ctx contextModel.AdminContext, data adminUserModel.UserAttrRequest) error {
	attr := adminUserModel.UserAttr{
		Model:     gorm.Model{ID: uint(data.Id)},
		SortId:    data.SortId,
		AttrType:  data.AttrType,
		AttrName:  data.AttrName,
		Color:     data.Color,
		HaveRange: data.HaveRange,
		Range:     data.Range,
	}

	attrs := make([]adminUserModel.UserAttr, 0)
	attrs = append(attrs, attr)

	err := adminUserRepo.SaveUserAttrs(attrs)
	if err != nil {
		log.Printf("adminUserService.UpdateAttr SaveUserAttrs err: %v", err)
		return err
	}

	return nil
}

func GetUserTags(ctx contextModel.AdminContext, page, limit int) ([]adminUserModel.UserTagResponse, int, error) {
	results := make([]adminUserModel.UserTagResponse, 0)

	tags, err := adminUserRepo.FindUserTags(page, limit)
	if err != nil {
		return results, 0, err
	}

	for _, tag := range tags {
		results = append(results, adminUserModel.UserTagResponse{
			Id:        int(tag.ID),
			SortId:    tag.SortId,
			Title:     tag.Title,
			GroupName: tag.GroupName,
		})
	}

	count, err := adminUserRepo.CountUserTag()
	if err != nil {
		return results, 0, err
	}

	return results, count, nil
}

func CreateUserTag(ctx contextModel.AdminContext, sortID int, title string, groupName string) error {
	tag := adminUserModel.UserTag{
		SortId:    sortID,
		Title:     title,
		GroupName: groupName,
	}

	err := adminUserRepo.SaveUserTag(&tag)
	return err
}

func UpdateUserTag(ctx contextModel.AdminContext, id, sortID int, title string) error {
	tag := adminUserModel.UserTag{
		Model:  gorm.Model{ID: uint(id)},
		SortId: sortID,
		Title:  title,
	}

	err := adminUserRepo.SaveUserTag(&tag)
	return err
}

func DeleteUserTag(ctx contextModel.AdminContext, id int) error {
	err := adminUserRepo.DeleteUserTag(id)
	return err
}
