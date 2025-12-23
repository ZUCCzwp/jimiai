package globalApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/globalModel"
	"jiyu/model/userModel"
	"jiyu/service/globalService"
	"jiyu/util"
	"jiyu/util/ishumei"
	"jiyu/util/response"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	ImagePath = "./runtime/upload/"
)

func UploadImages(c *gin.Context) {
	_ = MkUploadDir()

	form, err := c.MultipartForm()

	if err != nil {
		// 上传失败
		response.Error(c, response.ERROR, "上传失败", nil)
	}

	// 获取所有图片
	files := form.File["files"]

	result := make([]string, 0)

	for _, file := range files {

		ext := filepath.Ext(file.Filename)

		if !util.IsImageExt(ext) {
			response.Error(c, response.ERROR, "上传失败, 仅支持jpg, jpeg, png格式", nil)
			return
		}

		id, err := uuid.NewV4()

		if err != nil {
			log.Println("api.user.UploadAvatar uuid生成失败, error:", err)
			response.Error(c, response.ERROR, "上传失败", nil)
			return
		}

		filename := id.String() + ext

		dst := ImagePath + filename

		// 逐个存
		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Println("api.user.UploadAvatar 文件保存失败, error:", err)
			response.Error(c, response.ERROR, "上传失败", nil)
			return
		}

		url, err := globalService.UploadImage(dst)
		if err != nil {
			response.Error(c, response.ERROR, err.Error(), nil)
			return
		}

		result = append(result, url)
	}

	response.Success(c, "ok", result)
}

func ImgUploadToken(c *gin.Context) {
	var data globalModel.ImgUploadTokenRequest
	var category string
	err := c.ShouldBindJSON(&data)

	if err != nil {
		category = "upload"
	}

	ctx := c.MustGet("context").(contextModel.Context)

	token, err := globalService.ImgUploadToken(ctx, category, 1, 50*1024*1024)

	if err != nil {
		response.Error(c, response.ERROR, "获取上传token失败", nil)
	}

	response.Success(c, "ok", gin.H{
		"qnToken": token,
	})
}

func QiniuUploadToken(c *gin.Context) {
	category := c.Query("category")
	ctx := contextModel.Context{
		User: &userModel.User{},
	}
	if category == "" {
		category = "invite"
	}
	token, err := globalService.ImgUploadToken(ctx, category, 1, 50*1024*1024)
	if err != nil {
		response.Error(c, response.ERROR, "获取上传token失败", nil)
	}

	response.Success(c, "ok", gin.H{
		"qiniu_token": token,
	})
}

func MkUploadDir() error {
	if err := os.MkdirAll(ImagePath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Report 投诉
func Report(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	var data globalModel.ReportRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	err = globalService.Report(ctx, data)

	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func SMSCodeList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)

	pn := c.Query("page")
	page, err := strconv.ParseInt(pn, 10, 64)

	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	ps := c.Query("limit")
	limit, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		response.Error(c, -1, "参数错误", nil)
		return
	}

	list, count, err := globalService.SMSCodeList(ctx, int(page), int(limit))
	if err != nil {
		response.Error(c, -1, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}

func FindDoc(c *gin.Context) {
	// ctx := c.MustGet("context").(contextModel.AdminContext)

	doc, err := globalService.FindDoc()
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", doc)
}

func SaveDoc(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.AdminContext)
	var data globalModel.Doc
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	err = globalService.SaveDoc(ctx, &data)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", nil)
}

func DownloadRef(c *gin.Context) {
	id := c.Param("uid")
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.Error(c, response.ERROR, "参数错误", nil)
		return
	}

	ip := c.ClientIP()

	url, err := globalService.DownloadRef(int(uid), ip)
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func UserAttributesMap(c *gin.Context) {
	attrs, err := globalService.UserAttributesMap()
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", attrs)
}

func UserTagsMap(c *gin.Context) {
	tagsMap, err := globalService.UserTagsMap()
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", tagsMap)
}

func ShumeiText(c *gin.Context) {
	text := c.Query("text")

	b := ishumei.Text("13800000000", text, "")
	if b {
		response.Success(c, "ok", nil)
	} else {
		response.Error(c, response.ERROR, "error", nil)
	}
}

func FindArea(c *gin.Context) {
	areas, err := globalService.FindArea()
	if err != nil {
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", areas)
}
