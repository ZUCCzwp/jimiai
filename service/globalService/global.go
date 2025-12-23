package globalService

import (
	"errors"

	"jiyu/config"
	"jiyu/model/adminUserModel"
	"jiyu/model/contextModel"
	"jiyu/model/globalModel"
	"jiyu/repo/adminUserRepo"
	"jiyu/repo/globalRepo"
	"jiyu/repo/userRepo"
	"jiyu/util/qiniu"
	"log"
	"os"
	"time"

	"github.com/golang-module/carbon/v2"
)

func ImgUploadToken(ctx contextModel.Context, prefix string, fsizeMin, fsizeLimit int64) (string, error) {
	token, err := qiniu.GetImgUploadToken(ctx.User.ID, prefix, fsizeMin, fsizeLimit)
	if err != nil {
		log.Println("globalService.ImgUploadToken 七牛云上传token获取失败, error: ", err)
		return "", err
	}

	return token, err
}

func UploadImage(filePath string) (url string, err error) {
	filename, width, height, err := qiniu.PutFile(filePath)

	if err != nil {
		deleteLocalFile(filePath)
		log.Println("globalService.UploadImage 七牛云上传失败, error: ", err)
		return "", errors.New("上传失败")
	}

	image := globalModel.Image{
		UUID:   filename,
		Width:  int(width),
		Height: int(height),
	}

	err = globalRepo.CreateImage(image)

	if err != nil {
		deleteLocalFile(filePath)
		log.Println("globalService.UploadImage 图片数据存储失败, error: ", err)
		return "", errors.New("上传失败")
	}

	deleteLocalFile(filePath)
	return config.QiniuConfig.ImgHost + filename, nil
}

func deleteLocalFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Println("globalService.deleteLocalFile 删除失败, error: ", err)
	}
}

// Report 举报
func Report(ctx contextModel.Context, data globalModel.ReportRequest) error {

	report := data.ConvertToReport(ctx.User)

	err := globalRepo.CreateReport(report)

	if err != nil {
		log.Println("globalService.Report 举报失败, error: ", err)
		return errors.New("举报失败")
	}

	// 查看24小时内举报次数
	count, err := globalRepo.CountReportByUser(int(ctx.User.ID), data.TargetUid)

	if err != nil {
		log.Println("globalService.Report 查询举报次数失败, error: ", err)
		return errors.New("举报失败")
	}

	if count < 3 {
		return nil
	}

	if count >= 3 && count < 10 {
		// 封禁一天
		report.BanLevel = 1
		report.BanTime = carbon.Now().AddDay().ToStdTime()
		report.ProcessTime = time.Now()
		report.ProcessType = 0
	} else if count >= 10 && count < 20 {
		// 封禁7天
		report.BanLevel = 2
		report.BanTime = carbon.Now().AddWeek().ToStdTime()
		report.ProcessTime = time.Now()
		report.ReportType = 0
	} else if count >= 20 {
		// 封禁一个月
		report.BanLevel = 3
		report.BanTime = carbon.Now().AddMonth().ToStdTime()
		report.ProcessTime = time.Now()
		report.ReportType = 0
	}

	targetUser, err := userRepo.FindByID(data.TargetUid)
	if err != nil {
		log.Println("globalService.Report 查询用户失败, error: ", err)
		return errors.New("举报失败")
	}

	targetUser.BanLevel = report.BanLevel
	banTime := report.BanTime
	targetUser.BanTime = &banTime
	targetUser.BanRecordId = int(report.ID)

	err = userRepo.Save(targetUser)

	return nil
}

func SMSCodeList(ctx contextModel.AdminContext, page, limit int) ([]globalModel.SMSCodeResponse, int, error) {
	list, err := globalRepo.SMSCodeListForAdmin(page, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := globalRepo.SMSCodeCount()
	if err != nil {
		return nil, 0, err
	}

	for i := range list {
		list[i].MsgType = "验证码"
		list[i].SendType = "-"
	}

	return list, count, nil
}

func FindDoc() (*globalModel.Doc, error) {
	return globalRepo.FindDoc()
}

func SaveDoc(ctx contextModel.AdminContext, doc *globalModel.Doc) error {
	return globalRepo.UpdateDoc(doc)
}

func DownloadRef(refId int, ip string) (string, error) {
	// 记录 refID - 访客ip
	err := globalRepo.SaveRefIpMap(globalModel.RefIdIp{RefId: refId, Ip: ip})
	if err != nil {
		return "", err
	}

	// 获取下载地址 todo

	return "", nil
}

func UserAttributesMap() ([]adminUserModel.UserAttr, error) {
	return adminUserRepo.FindUserAttrs()
}

func UserTagsMap() (map[string][]adminUserModel.UserTagResponse, error) {
	tags, err := adminUserRepo.FindUserTags(1, 999)
	if err != nil {
		log.Println("globalService.UserTagsMap 查询用户标签失败, error: ", err)
		return nil, err
	}

	results := make(map[string][]adminUserModel.UserTagResponse)

	for _, v := range tags {
		if _, ok := results[v.GroupName]; !ok {
			results[v.GroupName] = make([]adminUserModel.UserTagResponse, 0)
		}

		results[v.GroupName] = append(results[v.GroupName], adminUserModel.UserTagResponse{
			Id:        int(v.ID),
			SortId:    v.SortId,
			Title:     v.Title,
			GroupName: v.GroupName,
		})
	}

	return results, nil
}

func FindArea() ([]globalModel.AreaResponse, error) {
	areas, err := globalRepo.FindArea()
	if err != nil {
		log.Println("查询地区数据失败:", err)
		return nil, err
	}
	// 将 areas 转换为 []globalModel.AreaResponse
	areaResponses := make([]globalModel.AreaResponse, 0)
	// 构建省市层级关系
	provinceMap := make(map[int]globalModel.AreaResponse)
	cityMap := make(map[int][]globalModel.AreaResponse)

	// 先将所有城市按parentId分组
	for _, area := range areas {
		if area.ParentID == 0 {
			// 省级
			provinceMap[area.AreaID] = globalModel.AreaResponse{
				Name:     area.Name,
				AreaID:   area.AreaID,
				Children: make([]globalModel.AreaResponse, 0),
			}
		} else {
			// 市级
			if _, ok := cityMap[area.ParentID]; !ok {
				cityMap[area.ParentID] = make([]globalModel.AreaResponse, 0)
			}
			cityMap[area.ParentID] = append(cityMap[area.ParentID], globalModel.AreaResponse{
				Name:     area.Name,
				AreaID:   area.AreaID,
				Children: nil,
			})
		}
	}

	// 组装省市数据
	for provinceId, province := range provinceMap {
		if cities, ok := cityMap[provinceId]; ok {
			province.Children = cities
			areaResponses = append(areaResponses, province)
		}
	}

	return areaResponses, nil
}
