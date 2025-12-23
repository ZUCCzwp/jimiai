package globalModel

import (
	"jiyu/model/userModel"
	"time"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	UUID   string `gorm:"type:varchar(255);not null;unique;index:idx_uuid"`
	Width  int    `gorm:"type:int;default:0"`
	Height int    `gorm:"type:int;default:0"`
}

type Report struct {
	gorm.Model
	UserId         int       `json:"user_id" gorm:"type:int;index:idx_user_id"`
	Nickname       string    `json:"nickname" gorm:"type:varchar(255)"` // 举报人昵称
	Phone          string    `json:"phone" gorm:"type:varchar(255);index:idx_phone"`
	TargetUid      int       `json:"target_uid" gorm:"type:int;index:idx_target_uid"`
	TargetNickname string    `json:"target_nickname" gorm:"type:varchar(255)"`
	Sign           string    `json:"sign" gorm:"type:varchar(255)"` // 个性签名
	Avatar         string    `json:"avatar" gorm:"type:varchar(255)"`
	BgImage        string    `json:"bg_image" gorm:"type:varchar(255)"`
	Content        string    `json:"content" gorm:"type:text"`               // 举报文字内容
	ReportScene    int       `json:"report_scene" gorm:"type:int"`           // 举报场景 0=漂流瓶 1=动态 2=消息 3=个人消息
	ReportType     int       `json:"report_type" gorm:"type:int"`            // 举报类型 0=低俗色情 1=垃圾广告 2=骚扰/不文明 3=涉嫌欺诈 4=政治 5=恐暴 6=其他
	Detail         string    `json:"detail" gorm:"type:text"`                // 用户写的举报补充说明
	Process        int       `json:"process" gorm:"type:int;default:0"`      // 0=未处理 1=已处罚 2=未违规
	ProcessTime    time.Time `json:"process_time" gorm:"type:datetime;"`     // 处理时间
	ProcessType    int       `json:"process_type" gorm:"type:int;default:0"` // 0=系统处理 1=人工处理
	BanLevel       int       `json:"ban_level" gorm:"type:int;default:0"`    // 0=不禁言 1=禁言一天 2=禁言一周 3=禁言一个月 4=永久禁言 5=封号
	BanTime        time.Time `json:"ban_time" gorm:"type:datetime;"`         // 封禁到期时间
}

type ImgUploadTokenRequest struct {
	Category string `json:"category"`
}

type ReportRequest struct {
	TargetUid  int      `json:"target_uid"`
	Nickname   string   `json:"nickname"`
	Sign       string   `json:"sign"` // 个性签名
	Avatar     string   `json:"avatar"`
	BgImage    string   `json:"bg_image"`
	Content    string   `json:"content"`     // 举报文字内容
	Images     []string `json:"images"`      // 举报的图片
	ReportType int      `json:"report_type"` // 举报类型 0=低俗色情 1=垃圾广告 2=骚扰/不文明 3=涉嫌欺诈 4=政治 5=恐暴 6=其他
	Detail     string   `json:"detail"`      // 用户写的举报补充说明
}

func (r *ReportRequest) ConvertToReport(user *userModel.User) Report {
	return Report{
		UserId:         int(user.ID),
		Phone:          user.Phone,
		TargetUid:      r.TargetUid,
		TargetNickname: r.Nickname,
		Sign:           r.Sign,
		Avatar:         r.Avatar,
		BgImage:        r.BgImage,
		Content:        r.Content,
		ReportType:     r.ReportType,
		Detail:         r.Detail,
		Process:        0,
		ProcessTime:    time.Now(),
		ProcessType:    0,
		BanLevel:       0,
	}
}

type Feedback struct {
	gorm.Model
	UserId        int       `json:"user_id" gorm:"type:int;index:idx_user_id"`
	Nickname      string    `json:"nickname" gorm:"type:varchar(255)"`        // 反馈人昵称
	SystemVersion string    `json:"system_version" gorm:"type:varchar(255)"`  // 系统版本
	PhoneModel    string    `json:"phone_model" gorm:"type:varchar(255)"`     // 手机型号
	FeedbackType  int       `json:"feedback_type" gorm:"type:int"`            // 反馈类型
	Content       string    `json:"content" gorm:"type:text"`                 // 反馈内容
	Status        int       `json:"status" gorm:"type:int;default:0"`         // 处理状态
	ProcessTime   time.Time `json:"process_time" gorm:"type:datetime;"`       // 处理时间
	Replay        string    `json:"replay" gorm:"type:text"`                  // 回复
	IsRead        bool      `json:"is_read" gorm:"type:tinyint(1);default:0"` // 是否已读
}

const (
	FeedbackStatusHandled   = 1
	FeedbackStatusUnhandled = 0
)

type FeedbackListResponse struct {
	ID           int      `json:"id"`
	Nickname     string   `json:"nickname"`
	FeedbackType int      `json:"feedback_type"`
	Content      string   `json:"content"`
	Images       []string `json:"images"`
	Status       int      `json:"status"`
	ProcessTime  string   `json:"process_time"`
	Replay       string   `json:"replay"`
	IsRead       bool     `json:"is_read"`
}

type QA struct {
	gorm.Model
	Question string `json:"question" gorm:"type:varchar(255)"`
	Answer   string `json:"answer" gorm:"type:text"`
}

type QAListResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type BannedListResponse struct {
	Id            int       `json:"id"`
	Uid           int       `json:"uid"`
	Nickname      string    `json:"nickname"`
	BanLevel      string    `json:"ban_level"`
	BanLevelIndex int       `json:"ban_level_index"`
	BanType       string    `json:"ban_type"`
	BanTime       time.Time `json:"ban_time"`
}

type SMSCode struct {
	gorm.Model
	Phone string `json:"phone" gorm:"type:varchar(255);index:idx_phone"`
	Code  string `json:"code" gorm:"type:varchar(255)"`
}

type EmailCode struct {
	gorm.Model
	Email string `json:"email" gorm:"type:varchar(255);index:idx_email"`
	Code  string `json:"code" gorm:"type:varchar(255)"`
}

type SMSCodeResponse struct {
	Id       int       `json:"id"`
	MsgType  string    `json:"msg_type"`
	Phone    string    `json:"phone"`
	Code     string    `json:"code"`
	Time     time.Time `json:"time"`
	SendType string    `json:"send_type"`
}

type Doc struct {
	gorm.Model
	// 用户协议
	UserAgreement string `json:"user_agreement" gorm:"type:mediumtext"`
	// 隐私政策
	PrivacyPolicy string `json:"privacy_policy" gorm:"type:mediumtext"`
	// 会员协议
	MemberAgreement string `json:"member_agreement" gorm:"type:mediumtext"`
}

type RefIdIp struct {
	gorm.Model
	RefId int    `json:"ref_id" gorm:"type:int;index:idx_ref_id"`
	Ip    string `json:"ip" gorm:"type:varchar(255);index:idx_ip"`
}

// 省市区
type Area struct {
	AreaID   int    `json:"area_id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
}

type AreaResponse struct {
	Name     string         `json:"name"`
	AreaID   int            `json:"area_id"`
	Children []AreaResponse `json:"children"`
}
