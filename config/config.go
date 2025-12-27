package config

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type ImageList []string

func (t *ImageList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t ImageList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

const (
	configFilePath = "./config/config.ini"
)

var (
	DBConfig          = Database{}
	AppConfig         = App{}
	AliyunDypnsConfig = AliyunDypns{}
	QiniuConfig       = Qiniu{}
	DomainConfig      = Domain{}
	RedisConfig       = redis{}
	GaodeConfig       = Gaode{}
	PaymentConfig     = PaymentInfo{}
	MemberConfig      = Member{}
	LocalConfig       = local{}
	HotUserConfig     = HotUser{}
	AdminAppConfig    = adminApp{}
	ShumeiConfig      = Shumei{}
	AdConfig          = Ad{}
	EmailConfig       = Email{}
)

type Database struct {
	DbHost      string
	DbPort      string
	UserName    string
	Password    string
	DBName      string
	TimeZone    string
	TablePrefix string
}

type App struct {
	DefaultAvatar           string `json:"default_avatar" gorm:"type:varchar(255)"`    // 默认头像地址
	DefaultDeleteImageCount int    `json:"default_delete_image_count" gorm:"type:int"` // 闪照每天销毁次数
	DyuAPIURL               string `json:"dyu_api_url" gorm:"type:varchar(255)"`       // API基础地址
	DyuAPIVersion           string `json:"dyu_api_version" gorm:"type:varchar(10)"`    // API版本号，如 v1, v2
}

// Domain 域名相关配置
type Domain struct {
	APIDomain string `json:"api_domain"` // 接口域名
	H5Domain  string `json:"h5_domain"`  // H5 页面域名
}

type AliyunDypns struct {
	AccessKeyId     string `json:"access_key_id" gorm:"type:varchar(255)"`     // 阿里云access key id
	AccessKeySecret string `json:"access_key_secret" gorm:"type:varchar(255)"` // 阿里云access key secret
}

type Qiniu struct {
	AccessKey string `json:"access_key" gorm:"type:varchar(255)"` // 七牛access key
	SecretKey string `json:"secret_key" gorm:"type:varchar(255)"` // 七牛secret key
	Bucket    string `json:"bucket" gorm:"type:varchar(255)"`     // 七牛bucket
	ImgHost   string `json:"img_host" gorm:"type:varchar(255)"`   // 七牛图片地址
}

type local struct {
	IsMaster bool   // 是否主服务器
	RunPort  string // 运行端口
}

type redis struct {
	Addr string
	Pwd  string
	DB   int
}

type HotUser struct {
	Count  int  `json:"count" gorm:"type:int"`         // 一次取多少个热门用户
	Enable bool `json:"enable" gorm:"type:tinyint(1)"` // 是否开启热门用户
}

type adminApp struct {
	RunPort string
}

// Gaode 高德地图
type Gaode struct {
	Key string `json:"key" gorm:"type:varchar(255)"` // 高德地图key
}

// Shumei 数美
type Shumei struct {
	AccessKey string `json:"access_key"`
	AppId     string `json:"app_id"`
}

// Ad 广告
type Ad struct {
	KaiPinAdEnable bool `json:"kai_pin_ad_enable" gorm:"type:tinyint(1)"` // 开屏广告是否开启
}

// Email 邮箱配置
type Email struct {
	SMTPHost     string `json:"smtp_host"`     // SMTP服务器地址，如：smtp.qq.com
	SMTPPort     string `json:"smtp_port"`     // SMTP端口，如：587
	FromEmail    string `json:"from_email"`    // 发送邮箱地址
	FromPassword string `json:"from_password"` // 发送邮箱密码或授权码
	FromName     string `json:"from_name"`     // 发送者名称
}

type Member struct {
	LevelRules string `json:"level_rules" gorm:"type:varchar(1024)"` // 等级规则
}

type PaymentInfo struct {
	AlipayInfo   AliPayment    `json:"alipay_config" gorm:"type:json"` // 支付宝支付开关
	WechatConfig WechatPayment `json:"wechat_config" gorm:"type:json"` // 微信支付开关
}

type AliPayment struct {
	AlipayEnable            bool   `json:"alipay_enable"`                                // 支付宝支付开关
	AlipayPartnerID         string `json:"alipay_partner_id" gorm:"type:text;"`          // 支付宝合作者身份ID
	AlipayLoginAccount      string `json:"alipay_seller_id" gorm:"type:text;"`           // 支付宝登陆账号
	AlipayAndroidPrivateKey string `json:"alipay_android_private_key" gorm:"type:text;"` // 支付宝安卓密钥
	AppId                   string `json:"app_id" gorm:"type:text;"`                     // 支付宝appid
	AlipayPublicKey         string `json:"alipay_public_key" gorm:"type:text;"`          // 支付宝公钥
	PrivateKey              string `json:"private_key" gorm:"type:text;"`                // 支付宝私钥
	IsProd                  bool   `json:"is_prod" gorm:"type:tinyint(1)"`               // 是否生产环境
	ReturnUrl               string `json:"return_url" gorm:"type:text;"`                 // 支付宝回调地址
	NotifyUrl               string `json:"notify_url" gorm:"type:text;"`                 // 支付宝回调地址
}

func (c AliPayment) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *AliPayment) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type WechatPayment struct {
	WechatEnable     bool   `json:"wechat_enable"`                        // 微信支付开关
	WechatAppID      string `json:"wechat_app_id" gorm:"type:text;"`      // 微信开放平台移动应用AppID
	WechatAppSecret  string `json:"wechat_app_secret" gorm:"type:text;"`  // 微信开放平台移动应用AppSecret
	WechatMchID      string `json:"wechat_mch_id" gorm:"type:text;"`      // 微信开放平台绑定商户号mch id
	WechatSerialNo   string `json:"wechat_serial_no" gorm:"type:text;"`   // 微信开放平台绑定商户号序列号
	WechatPrivateKey string `json:"wechat_private_key" gorm:"type:text;"` // 微信开放平台绑定商户号密钥key
	WechatAPIv3Key   string `json:"wechat_api_v3_key" gorm:"type:text;"`  // 微信开放平台绑定商户号APIv3密钥key
	WechatNotifyUrl  string `json:"wechat_notify_url" gorm:"type:text;"`  // 微信开放平台绑定商户号回调地址
}

func (c WechatPayment) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *WechatPayment) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type PaymentResponse struct {
	Enable bool   `json:"enable"` // 支付开关
	Type   string `json:"type"`   // 支付类型
}

func loadConfigSection(cfg *ini.File, section string, target interface{}) error {
	err := cfg.Section(section).MapTo(target)
	if err != nil {
		log.Fatalln("配置文件解析失败，请检查配置文件路径是否正确")
	}
	return err
}

func loadRequiredConfigs(cfg *ini.File) error {
	requiredSections := map[string]interface{}{
		"Database": &DBConfig,
		"Redis":    &RedisConfig,
		"Local":    &LocalConfig,
		"AdminApp": &AdminAppConfig,
		"Domain":   &DomainConfig,
	}

	for section, target := range requiredSections {
		if err := loadConfigSection(cfg, section, target); err != nil {
			return err
		}
	}
	return nil
}

// expandEnv 替换字符串中的环境变量
// 支持 ${VAR} 和 ${VAR:-default} 格式
func expandEnv(s string) string {
	return os.Expand(s, func(key string) string {
		// 处理 ${VAR:-default} 格式
		if idx := strings.Index(key, ":-"); idx > 0 {
			envKey := key[:idx]
			defaultValue := key[idx+2:]
			if value := os.Getenv(envKey); value != "" {
				return value
			}
			return defaultValue
		}
		// 处理 ${VAR} 格式
		return os.Getenv(key)
	})
}

// expandEnvInConfig 遍历配置文件中所有值并替换环境变量
func expandEnvInConfig(cfg *ini.File) {
	for _, section := range cfg.Sections() {
		for _, key := range section.Keys() {
			value := key.Value()
			expanded := expandEnv(value)
			if expanded != value {
				key.SetValue(expanded)
			}
		}
	}
}

func InitConfig() {
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalln("配置文件读取失败，请检查配置文件路径是否正确")
	}

	// 替换配置文件中的环境变量
	expandEnvInConfig(cfg)

	if err := loadRequiredConfigs(cfg); err != nil {
		log.Fatalln(err)
	}

	// 加载可选配置
	if cfg.HasSection("Email") {
		if err := loadConfigSection(cfg, "Email", &EmailConfig); err != nil {
			log.Printf("邮箱配置加载失败: %v", err)
		}
	}

	if cfg.HasSection("App") {
		if err := loadConfigSection(cfg, "App", &AppConfig); err != nil {
			log.Printf("App配置加载失败: %v", err)
		}
	}
}

func LoadConfig(configFilePath string) {
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalln("配置文件读取失败，请检查配置文件路径是否正确")
	}

	// 替换配置文件中的环境变量
	expandEnvInConfig(cfg)

	if err := loadRequiredConfigs(cfg); err != nil {
		log.Fatalln(err)
	}

	// 加载可选配置
	if cfg.HasSection("Email") {
		if err := loadConfigSection(cfg, "Email", &EmailConfig); err != nil {
			log.Printf("邮箱配置加载失败: %v", err)
		}
	}

	if cfg.HasSection("App") {
		if err := loadConfigSection(cfg, "App", &AppConfig); err != nil {
			log.Printf("App配置加载失败: %v", err)
		}
	}
}
