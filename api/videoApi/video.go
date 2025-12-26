package videoApi

import (
	"jiyu/model/contextModel"
	"jiyu/model/videoModel"
	"jiyu/service/videoService"
	"jiyu/util/response"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// GenerateVideo 图生视频接口
func GenerateVideo(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("input_reference")
	if err != nil {
		log.Printf("videoApi.GenerateVideo 文件获取失败, error: %v", err)
		response.Error(c, response.INVALID_PARAMS, "请上传图片文件", nil)
		return
	}

	// 获取其他表单参数
	prompt := c.PostForm("prompt")
	model := c.PostForm("model")
	size := c.PostForm("size")
	seconds := c.PostForm("seconds")

	// 验证必填参数
	if prompt == "" || model == "" || size == "" || seconds == "" {
		log.Printf("videoApi.GenerateVideo 参数不完整")
		response.Error(c, response.INVALID_PARAMS, "参数不完整，请提供 prompt、model、size、seconds", nil)
		return
	}

	// 从 context 获取用户信息
	ctx := c.MustGet("context").(contextModel.Context)
	user := ctx.User

	// 从用户表获取 API Key
	apiKey := user.ApiKey
	if apiKey == "" {
		log.Printf("videoApi.GenerateVideo 用户 API Key 未配置, userId: %d", user.ID)
		response.Error(c, response.ERROR, "API Key 未配置，请先设置您的 API Key", nil)
		return
	}

	// 创建临时文件保存上传的图片
	tempDir := "./runtime/upload/"
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		log.Printf("videoApi.GenerateVideo 创建临时目录失败, error: %v", err)
		response.Error(c, response.ERROR, "服务器错误", nil)
		return
	}

	// 保存文件到临时目录
	tempFilePath := tempDir + file.Filename
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		log.Printf("videoApi.GenerateVideo 文件保存失败, error: %v", err)
		response.Error(c, response.ERROR, "文件保存失败", nil)
		return
	}

	// 确保临时文件在使用后删除
	defer func() {
		if err := os.Remove(tempFilePath); err != nil {
			log.Printf("videoApi.GenerateVideo 删除临时文件失败, error: %v", err)
		}
	}()

	// 构建请求
	req := videoModel.VideoGenerationRequest{
		Prompt:  prompt,
		Model:   model,
		Size:    size,
		Seconds: seconds,
	}

	// 调用 service 层
	result, err := videoService.GenerateVideo(req, tempFilePath, apiKey)
	if err != nil {
		log.Printf("videoApi.GenerateVideo 调用服务失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	// 直接返回外部API的响应，保持格式一致
	response.Success(c, "ok", result)
}

// GetVideo 查询视频任务接口
func GetVideo(c *gin.Context) {
	// 从 URL 参数获取视频ID
	videoID := c.Param("videoId")
	if videoID == "" {
		log.Printf("videoApi.GetVideo 视频ID不能为空")
		response.Error(c, response.INVALID_PARAMS, "视频ID不能为空", nil)
		return
	}

	// 从 context 获取用户信息
	ctx := c.MustGet("context").(contextModel.Context)
	user := ctx.User

	// 从用户表获取 API Key
	apiKey := user.ApiKey
	if apiKey == "" {
		log.Printf("videoApi.GetVideo 用户 API Key 未配置, userId: %d", user.ID)
		response.Error(c, response.ERROR, "API Key 未配置，请先设置您的 API Key", nil)
		return
	}

	// 调用 service 层
	result, err := videoService.GetVideo(videoID, apiKey)
	if err != nil {
		log.Printf("videoApi.GetVideo 调用服务失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	// 直接返回外部API的响应，保持格式一致
	response.Success(c, "ok", result)
}

