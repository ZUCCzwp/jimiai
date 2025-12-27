package videoService

import (
	"jiyu/config"
	"jiyu/model/videoModel"
	"jiyu/util/http"
	"log"
	"strings"
)

// GenerateVideo 调用图生视频接口
func GenerateVideo(req videoModel.VideoGenerationRequest, filePath string, apiKey string) (*videoModel.VideoGenerationResponse, error) {
	// 验证 API Key
	if apiKey == "" {
		log.Println("videoService.GenerateVideo API Key 未配置")
		return nil, &ServiceError{Message: "API Key 未配置"}
	}

	// 从配置文件获取 API 基础 URL 和版本号，必须配置
	baseURL := config.AppConfig.DyuAPIURL
	apiVersion := config.AppConfig.DyuAPIVersion
	if baseURL == "" {
		log.Println("videoService.GenerateVideo API地址未配置")
		return nil, &ServiceError{Message: "API地址未配置，请在配置文件中设置 dyu_api_url"}
	}
	if apiVersion == "" {
		apiVersion = "v1" // 默认使用 v1
	}

	// 拼接完整的 API URL
	dyuAPIURL := strings.TrimSuffix(baseURL, "/") + "/" + apiVersion + "/videos"

	// 构建表单字段
	formFields := map[string]string{
		"prompt":  req.Prompt,
		"model":   req.Model,
		"size":    req.Size,
		"seconds": req.Seconds,
	}

	// 调用外部 API
	var response videoModel.VideoGenerationResponse
	status, body, err := http.PostMultipart(dyuAPIURL, filePath, formFields, apiKey, &response)

	if err != nil {
		log.Printf("videoService.GenerateVideo 调用外部API失败, error: %v", err)
		return nil, err
	}

	// 即使状态码不是200，也返回解析后的响应体（外部API可能返回错误信息，但格式一致）
	if status != 200 {
		log.Printf("videoService.GenerateVideo 外部API返回错误状态码: %d, body: %s", status, body)
		// 如果解析失败，尝试手动解析
		if response.ID == "" {
			// 如果解析失败，返回错误
			return nil, &ServiceError{Message: "调用外部API失败: " + body}
		}
	}

	return &response, nil
}

// GetVideo 查询视频任务状态
func GetVideo(videoID string, apiKey string) (*videoModel.VideoQueryResponse, error) {
	// 验证 API Key
	if apiKey == "" {
		log.Println("videoService.GetVideo API Key 未配置")
		return nil, &ServiceError{Message: "API Key 未配置"}
	}

	// 从配置文件获取 API 基础 URL 和版本号，必须配置
	baseURL := config.AppConfig.DyuAPIURL
	apiVersion := config.AppConfig.DyuAPIVersion
	if baseURL == "" {
		log.Println("videoService.GetVideo API地址未配置")
		return nil, &ServiceError{Message: "API地址未配置，请在配置文件中设置 dyu_api_url"}
	}
	if apiVersion == "" {
		apiVersion = "v1" // 默认使用 v1
	}

	// 拼接完整的 API URL（包含视频ID）
	dyuAPIURL := strings.TrimSuffix(baseURL, "/") + "/" + apiVersion + "/videos/" + videoID

	// 调用外部 API
	var response videoModel.VideoQueryResponse
	status, body, err := http.Get(dyuAPIURL, apiKey, &response)

	if err != nil {
		log.Printf("videoService.GetVideo 调用外部API失败, error: %v", err)
		return nil, err
	}

	// 即使状态码不是200，也返回解析后的响应体（外部API可能返回错误信息，但格式一致）
	if status != 200 {
		log.Printf("videoService.GetVideo 外部API返回错误状态码: %d, body: %s", status, body)
		// 如果解析失败，尝试手动解析
		if response.ID == "" {
			// 如果解析失败，返回错误
			return nil, &ServiceError{Message: "调用外部API失败: " + body}
		}
	}

	return &response, nil
}

// ServiceError 服务错误
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
