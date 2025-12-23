package userModel

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"jiyu/config"
	"jiyu/model/memberModel"
	"jiyu/model/positionModel"
	"time"

	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserTags []int

// Value 实现方法
func (p UserTags) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *UserTags) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type User struct {
	gorm.Model
	Phone           string                 `json:"phone" gorm:"type:varchar(20);not null;unique"`
	Email           string                 `json:"email" gorm:"type:varchar(255);not null;unique"`
	Nickname        string                 `json:"nickname" gorm:"type:varchar(20);comment:'昵称'"`
	Password        string                 `json:"password" gorm:"type:varchar(255);not null;default:'';comment:'密码'"`
	Avatar          string                 `json:"avatar_url" gorm:"type:varchar(255);default:'';comment:'头像'"`
	City            string                 `json:"city" gorm:"type:varchar(255);default:'';comment:'城市'	"`
	Birthday        string                 `json:"birthday" gorm:"type:varchar(255);default:'';comment:'出生日期'"`
	Weight          int                    `json:"weight" gorm:"type:int;default:0;comment:'体重'"`
	SumVipDays      int                    `json:"sum_vip_days" gorm:"type:int;default:0;comment:'累积会员天数'"`                 // 累积会员天数
	HotBanned       bool                   `json:"hot_banned" gorm:"type:tinyint;default:0;comment:'是否被热度封禁'"`              // 是否被热度封禁
	BanLevel        int                    `json:"ban_level" gorm:"type:int;default:0;comment:'封禁等级'"`                      // 封禁等级 0=不封禁 1=禁言一天 2=禁言一周 3=禁言一个月 4=禁言1年 5=永久禁言 6=永久封号
	BanRecordId     int                    `json:"ban_record_id" gorm:"type:int;default:0"`                                 // 封禁记录id
	LastLoginIp     string                 `json:"last_login_ip" gorm:"type:varchar(20);default:'';comment:'最后一次登录ip'"`     // 最后一次登录ip
	RegisterDevice  string                 `json:"register_device" gorm:"type:varchar(255);default:'';comment:'注册设备'"`      // 注册设备
	RefId           int                    `json:"ref_id" gorm:"type:int;default:0;comment:'推荐人id'"`                        // 推荐人id
	RefCode         string                 `json:"ref_code" gorm:"type:varchar(255);not null;unique;comment:'推荐码'"`         // 推荐码
	MyInviteCode    string                 `json:"my_invite_code" gorm:"type:varchar(255);not null;unique;comment:'自身邀请码'"` // 自身邀请码
	VipTime         *time.Time             `json:"vip_time" gorm:"type:datetime;default:null"`                              // vip到期时间
	BanTime         *time.Time             `json:"ban_time" gorm:"type:datetime;default:null"`                              // 封禁到期时间
	LastRequestTime time.Time              `json:"last_request_time" gorm:"type:datetime;"`                                 // 最后一次请求时间
	Tags            UserTags               `json:"tags" gorm:"type:varchar(255)"`                                           // 用户标签
	Position        positionModel.Position `json:"position" gorm:"embedded;"`
	JimiCoin        int                    `json:"jimicoin" gorm:"type:int;default:0;comment:'jimi币'"`          // jimi币
	ApiKey          string                 `json:"api_key" gorm:"type:varchar(255);default:'';comment:'API密钥'"` // API密钥
}

func New(email string) *User {
	return &User{
		Email: email,
		Tags:  make(UserTags, 0),
		Position: positionModel.Position{
			Longitude: 0,
			Latitude:  0,
		},
		City:            "",
		LastRequestTime: time.Now(),
		HotBanned:       false,
		BanLevel:        0,
		BanRecordId:     0,
		LastLoginIp:     "",
		RegisterDevice:  "Android",
		JimiCoin:        0,
	}
}

func (u *User) IsBanned() (bool, error) {
	if u.BanTime != nil && u.BanTime.After(time.Now()) {
		return true, errors.New("您已被禁言, 截至时间: " + u.BanTime.Format("2006-01-02 15:04:05"))
	}
	return false, nil
}

func (u *User) IsVip() bool {
	return u.VipTime != nil && u.VipTime.After(time.Now())
}

func (u *User) GenerateRandomAvatar() string {
	avatar := config.AppConfig.DefaultAvatar
	return avatar
}

func (u *User) VipLevel(memberList []memberModel.MemberConfig) int {
	// 根据type建立映射
	memberMap := lo.SliceToMap(memberList, func(v memberModel.MemberConfig) (int, memberModel.MemberConfig) {
		return v.Type, v
	})
	// 根据累积会员天数计算会员等级
	// 当前累积天数在哪个范围那就是什么等级
	switch {
	case memberMap[memberModel.MemberTypeDiamond].MinDays <= u.SumVipDays:
		return memberModel.MemberTypeDiamond // 钻石会员
	case memberMap[memberModel.MemberTypeGolden].MinDays <= u.SumVipDays &&
		u.SumVipDays <= memberMap[memberModel.MemberTypeGolden].MaxDays:
		return memberModel.MemberTypeGolden // 黄金会员
	case memberMap[memberModel.MemberTypeSilver].MinDays <= u.SumVipDays &&
		u.SumVipDays <= memberMap[memberModel.MemberTypeSilver].MaxDays:
		return memberModel.MemberTypeSilver // 白银会员
	case memberMap[memberModel.MemberTypeBronze].MinDays <= u.SumVipDays &&
		u.SumVipDays <= memberMap[memberModel.MemberTypeBronze].MaxDays:
		return memberModel.MemberTypeBronze // 青铜会员
	default:
		return 0 // 普通用户
	}
}

// MaxDeleteImageCount 每天最大销毁照片数量
func (u *User) MaxDeleteImageCount() int {
	if u.IsVip() {
		return 999999
	}

	return config.AppConfig.DefaultDeleteImageCount
}

// Ban 封禁
// banLevel 1=禁言一天 2=禁言一周 3=禁言一个月 4=禁言1年 5=永久禁言 6=永久封号
func (u *User) Ban(banLevel int) time.Time {

	var level, day int
	switch banLevel {
	case 0:
		level = 0
		day = 0
	case 1:
		level = 1
		day = 1
	case 2:
		level = 2
		day = 7
	case 3:
		level = 3
		day = 30
	case 4:
		level = 4
		day = 365
	case 5:
		level = 5
		day = 99999
	case 6:
		level = 6
		day = 99999
	default:
		return time.Time{}
	}

	u.BanLevel = level
	if day > 0 {
		banTime := time.Now().Add(time.Hour * 24 * time.Duration(day))
		u.BanTime = &banTime
		u.BanRecordId = -1
		return *u.BanTime
	} else {
		u.BanTime = nil
		u.BanRecordId = 0
		return time.Time{}
	}
}

func (u *User) GetVipInfo(memberList []memberModel.MemberConfig) VipInfo {
	vipTime := time.Time{}
	if u.VipTime != nil {
		vipTime = *u.VipTime
	}
	return VipInfo{
		IsVip:      u.IsVip(),
		VipTime:    vipTime,
		VipLevel:   u.VipLevel(memberList),
		SumVipDays: u.SumVipDays,
	}
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// ComparePassword 比较密码
func (u *User) ComparePassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	return err == nil
}

type UserNickName struct {
	Nickname string `json:"nickname" gorm:"type:varchar(20);"`
}

type VipInfo struct {
	VipTime    time.Time `json:"vip_time"`     //会员到期时间
	VipLevel   int       `json:"vip_level"`    //会员等级
	SumVipDays int       `json:"sum_vip_days"` //累积会员天数
	IsVip      bool      `json:"is_vip"`       //是否会员
}

type LoginRequest struct {
	Nickname string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type EmailCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type AutoLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

type RegisterRequest struct {
	Phone      string `json:"phone"`
	Email      string `json:"email" binding:"required"`
	Nickname   string `json:"nickname" binding:"required"`
	Password   string `json:"password" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
	InviteCode string `json:"invite_code"` //邀请码
}

type UserAvatarRequest struct {
	Avatar string `form:"avatar" json:"avatar" binding:"required"`
	Width  int    `form:"width" json:"width"`
	Height int    `form:"height" json:"height"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type InfoResponse struct {
	Id           int             `json:"id"`
	Phone        string          `json:"phone"`
	Email        string          `json:"email"`
	Nickname     string          `json:"nickname"`
	Age          int             `json:"age"`
	MyInviteCode string          `json:"my_invite_code"`
	City         string          `json:"city"`
	VipInfo      VipInfo         `json:"vip_info"`
	Payment      PaymentResponse `json:"payment"`
	ApiKey       string          `json:"api_key"`
}

type UserLikeResponse struct {
	Shape int `json:"shape" gorm:"type:int;default:-1"` //0=狒狒 1=猴子 2=匀称 3=肌肉 4=熊猪 5=其他
}

type PaymentResponse struct {
	JimiCoin     int `json:"jimicoin" gorm:"type:int;default:0;"`      // 吉米币
	UsedCoin     int `json:"used_coin" gorm:"type:int;default:0;"`     // 已使用币
	RequestCount int `json:"request_count" gorm:"type:int;default:0;"` // 请求次数
}

// Detail 用户的详情
type Detail struct {
	Uid            int       `json:"uid"`
	Phone          string    `json:"phone"`
	BgImgUrl       string    `json:"bg_img_url"`
	City           string    `json:"city"`
	Distance       float64   `json:"distance"`
	LastOnlineTime time.Time `json:"last_online_time"`
	Tags           UserTags  `json:"tags"`
	MatchRate      int       `json:"match_rate"` // 匹配率
	FaceId         bool      `json:"face_id"`    // 是否视频认证
}

type IsBanResponse struct {
	IsBan   bool      `json:"is_ban"`
	BanTime time.Time `json:"ban_time"`
}

type HomeDataResponse struct {
	Id            int       `json:"id"`
	City          string    `json:"city"`
	Distance      float64   `json:"distance"`
	LastLoginTime time.Time `json:"last_login_time"`
	VipInfo       VipInfo   `json:"vip_info"` // 会员等级
}
