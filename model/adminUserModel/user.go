package adminUserModel

import (
	"jiyu/config"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(20);not null;unique;index:idx_username" json:"username"`
	Password     string `gorm:"type:varchar(255);not null" json:"password"`
	Role         int    `gorm:"not null;default:1" json:"role"` // 0=管理员 1=超级管理员
	Avatar       string `gorm:"type:varchar(255);not null;default:''" json:"avatar"`
	Name         string `gorm:"type:varchar(20);not null;default:''" json:"name"`
	Introduction string `gorm:"type:varchar(255);not null;default:''" json:"introduction"`
}

type AddUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int    `json:"role"`
}

type ResetPwdRequest struct {
	Id       int    `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetAuthorityRequest struct {
	Id   int `json:"id" binding:"required"`
	Role int `json:"role"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (u *User) ComparePwd(pwd string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	return !(err != nil)
}

func (u *User) HashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
}

func NewUser() *User {
	u := &User{
		Username:     "admin",
		Password:     "admin666",
		Role:         1,
		Avatar:       "http://cdn1.jiyujiaoyou.cn/admin/robot.gif",
		Name:         "超级管理员",
		Introduction: "超级管理员",
	}

	u.HashPassword()

	return u
}

func (u *User) TableName() string {
	return config.DBConfig.TablePrefix + "admin_users"
}

type ChangePwdRequest struct {
	OldPwd string `json:"old_pwd" binding:"required"`
	NewPwd string `json:"new_pwd" binding:"required"`
}

type InfoResponse struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}

type UserListResponse struct {
	Id       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Birthday string `json:"birthday"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	RegisterDevice string    `json:"register_device"`
	RegisterTime   time.Time `json:"register_time"`
	LastLoginTime  time.Time `json:"last_login_time"`
	LastLoginIp    string    `json:"last_login_ip"`
	Status         int       `json:"status"` // 0=正常 1=禁言 2=封号
}

type UserBanRequest struct {
	UserIDRequest
	BanLevel int `json:"ban_level"`
}

type UserIDRequest struct {
	Uid int `json:"uid"`
}

type UserDetailResponse struct {
	Id    int    `json:"id"`
	Phone string `json:"phone"`
}

type UpdateUserRequest struct {
	UserDetailResponse
}

type RefListResponse struct {
	Id          int       `json:"id"`
	Uid         int       `json:"uid"`
	Nickname    string    `json:"nickname"`
	RefId       int       `json:"ref_id"`
	RefNickname string    `json:"ref_nickname"` // 邀请人昵称
	Time        time.Time `json:"time"`
}

type RefIncomeResponse struct {
	Id       int     `json:"id"`
	Uid      int     `json:"uid"`
	Nickname string  `json:"nickname"`
	Income   float64 `json:"income"`
	Count    int     `json:"count"`
}
