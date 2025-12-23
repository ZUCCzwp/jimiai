package settingModel

import (
	"database/sql/driver"
	"encoding/json"
	"jiyu/config"

	"gorm.io/gorm"
)

type IntList []int

func (t *IntList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t IntList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type LoginSetting struct {
	EnableSMS bool   `json:"enable_sms"`                    // 短信验证码开关
	SMSApiURL string `json:"sms_api_url" gorm:"type:text;"` // 短信验证码接口地址

	QcloudSMSAppID  string `json:"qcloud_sms_app_id" gorm:"type:text;"`  // 腾讯云短信AppID
	QcloudSMSAppKey string `json:"qcloud_sms_app_key" gorm:"type:text;"` // 腾讯云短信AppKey

	QcloudSMSSign       string `json:"qcloud_sms_sign" gorm:"type:text;"`        // 腾讯云短信签名
	QcloudSMSTemplateID string `json:"qcloud_sms_template_id" gorm:"type:text;"` // 腾讯云短信模板ID

	QcloudSMSForeignSign       string `json:"qcloud_sms_foreign_sign" gorm:"type:text;"`        // 腾讯云短信国外签名
	QcloudSMSForeignTemplateID string `json:"qcloud_sms_foreign_template_id" gorm:"type:text;"` // 腾讯云短信国外模板ID

	EnableSMSIPLimit bool `json:"enable_sms_ip_limit"` // 短信验证码IP限制开关
	SMSIPLimitCount  int  `json:"sms_ip_limit_count"`  // 短信验证码IP限制次数
}

type WithdrawSetting struct {
	MinWithdrawAmount float64 `json:"min_withdraw_amount" gorm:"type:decimal(9,2)"` // 最小提现金额
	WithdrawRate      float64 `json:"withdraw_rate" gorm:"type:decimal(9,2)"`       // 提现比例
}

type IMSetting struct {
	IMSdkAppID  string `json:"im_sdk_app_id" gorm:"type:text;"`  // 即时通讯SDKAppID
	ImSdkAppKey string `json:"im_sdk_app_key" gorm:"type:text;"` // 即时通讯SDKAppKey

	AndroidAccessID  string `json:"android_access_id" gorm:"type:text;"`  // Android推送AccessID
	AndroidSecretKey string `json:"android_secret_key" gorm:"type:text;"` // Android推送AccessKey

	IOSAccessID  string `json:"ios_access_id" gorm:"type:text;"`  // IOS推送AccessID
	IOSSecretKey string `json:"ios_secret_key" gorm:"type:text;"` // IOS推送AccessKey
}

type InviteSetting struct {
	EnableInvite      bool    `json:"enable_invite"`                                // 邀请开关
	InviteRewardValue float64 `json:"invite_reward_value" gorm:"type:decimal(9,2)"` // 邀请奖励币值
}

type AppSetting struct {
	config.App
}

type ReportSetting struct {
	ReportTime        IntList `json:"report_time" gorm:"type:json"` // 举报次数(超过就处罚)
	FailedAuditsCount int     `json:"failed_audits_count"`          // 审核不通过次数
}

type AliyunSetting struct {
	config.AliyunDypns
}

type QiniuSetting struct {
	config.Qiniu
}

type GaodeSetting struct {
	config.Gaode
}

type HotUserSetting struct {
	config.HotUser
}

type ShumeiSetting struct {
	config.Shumei
}

type AdSetting struct {
	config.Ad
}

type PaymentSetting struct {
	config.PaymentInfo
}

type MemberSetting struct {
	config.Member
}

type Setting struct {
	gorm.Model
	Login    LoginSetting    `gorm:"embedded" json:"login"`
	Withdraw WithdrawSetting `gorm:"embedded" json:"withdraw"`
	IM       IMSetting       `gorm:"embedded" json:"im"`
	Payment  PaymentSetting  `gorm:"embedded" json:"payment"`
	Invite   InviteSetting   `gorm:"embedded" json:"invite"`
	App      AppSetting      `gorm:"embedded" json:"app"`
	Report   ReportSetting   `gorm:"embedded" json:"report"`
	Aliyun   AliyunSetting   `gorm:"embedded;embeddedPrefix:aliyun_" json:"aliyun"`
	Qiniu    QiniuSetting    `gorm:"embedded;embeddedPrefix:qiniu_" json:"qiniu"`
	Gaode    GaodeSetting    `gorm:"embedded;embeddedPrefix:gaode_" json:"gaode"`
	HotUser  HotUserSetting  `gorm:"embedded;embeddedPrefix:hot_user_" json:"hot_user"`
	Shumei   ShumeiSetting   `gorm:"embedded;embeddedPrefix:shumei_" json:"shumei"`
	Ad       AdSetting       `gorm:"embedded;embeddedPrefix:ad_" json:"ad"`
	Member   MemberSetting   `gorm:"embedded;embeddedPrefix:member_" json:"member"`
}

func (s *Setting) InitToConfig() {
	config.HotUserConfig = s.HotUser.HotUser
	config.GaodeConfig = s.Gaode.Gaode
	config.QiniuConfig = s.Qiniu.Qiniu
	config.AliyunDypnsConfig = s.Aliyun.AliyunDypns
	config.ShumeiConfig = s.Shumei.Shumei
	config.AdConfig = s.Ad.Ad
	config.PaymentConfig = s.Payment.PaymentInfo
	config.MemberConfig = s.Member.Member
}
