package globalRepo

import (
	"jiyu/global"
	"jiyu/model/globalModel"
)

func CreateImage(image globalModel.Image) error {
	return global.DB.Model(&globalModel.Image{}).Create(&image).Error
}

func FindImage(filename string) (*globalModel.Image, error) {
	var image globalModel.Image
	err := global.DB.Model(&globalModel.Image{}).Where("uuid = ?", filename).First(&image).Error
	return &image, err
}

// CreateReport 创建举报记录
func CreateReport(report globalModel.Report) error {
	return global.DB.Model(&globalModel.Report{}).Create(&report).Error
}

// CountReportByUser 查询24小时内举报次数
func CountReportByUser(uid int, targetUid int) (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.Report{}).Where("user_id = ? AND target_uid = ? AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR)", uid, targetUid).Count(&count).Error
	return int(count), err
}

func SMSCodeListForAdmin(page, limit int) ([]globalModel.SMSCodeResponse, error) {
	var results []globalModel.SMSCodeResponse
	err := global.DB.Model(&globalModel.SMSCode{}).Select("phone, code, created_at as time, id").Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&results).Error

	// 记录到redis
	// _, err = redis.RDB.Set(redis.Ctx, redis.KeySMSCode+fmt.Sprintf("::%s", phone), strconv.Itoa(code), time.Minute*10).Result()
	// if err != nil {
	// 	log.Println("userService.LoginSMSCode redis写入失败, error: ", err)
	// 	return "", errors.New("发送验证码失败")
	// }
	return results, err
}

func SMSCodeCount() (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.SMSCode{}).Count(&count).Error
	return int(count), err
}

func DocCount() (int, error) {
	var count int64
	err := global.DB.Model(&globalModel.Doc{}).Count(&count).Error
	return int(count), err
}

func SaveDoc(doc *globalModel.Doc) error {
	return global.DB.Save(doc).Error
}

func FindDoc() (*globalModel.Doc, error) {
	var result globalModel.Doc
	err := global.DB.Model(&globalModel.Doc{}).Where("id = 1").First(&result).Error
	return &result, err
}

func UpdateDoc(doc *globalModel.Doc) error {
	m := make(map[string]interface{})
	m["user_agreement"] = doc.UserAgreement
	m["privacy_policy"] = doc.PrivacyPolicy
	m["member_agreement"] = doc.MemberAgreement
	err := global.DB.Model(&globalModel.Doc{}).Where("id = 1").Updates(m).Error
	return err
}

// SaveRefIpMap 保存邀请人和ip的映射关系
func SaveRefIpMap(data globalModel.RefIdIp) error {
	return global.DB.Model(&globalModel.RefIdIp{}).Create(&data).Error
}

func FindRefIdByIp(ip string) (int, error) {
	var result globalModel.RefIdIp
	err := global.DB.Model(&globalModel.RefIdIp{}).Where("ip = ?", ip).First(&result).Error
	return result.RefId, err
}

func FindArea() ([]globalModel.Area, error) {
	var results []globalModel.Area
	err := global.DB.Model(&globalModel.Area{}).Find(&results).Error
	return results, err
}
