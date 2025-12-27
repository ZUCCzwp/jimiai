package watermarkService

import (
	"jiyu/config"
	"jiyu/model/watermarkModel"
	"jiyu/util/http"
	"log"
	"strings"
)

// RemoveWatermark 调用去水印接口
func RemoveWatermark(req watermarkModel.ChatCompletionRequest, apiKey string) (*watermarkModel.ChatCompletionResponse, error) {
	// 验证 API Key
	if apiKey == "" {
		log.Println("watermarkService.RemoveWatermark API Key 未配置")
		return nil, &ServiceError{Message: "API Key 未配置"}
	}

	// 从配置文件获取 API 基础 URL 和版本号，必须配置
	baseURL := config.AppConfig.DyuAPIURL
	apiVersion := config.AppConfig.DyuAPIVersion
	if baseURL == "" {
		log.Println("watermarkService.RemoveWatermark API地址未配置")
		return nil, &ServiceError{Message: "API地址未配置，请在配置文件中设置 dyu_api_url"}
	}
	if apiVersion == "" {
		apiVersion = "v1" // 默认使用 v1
	}

	// 拼接完整的 API URL
	dyuAPIURL := strings.TrimSuffix(baseURL, "/") + "/" + apiVersion + "/chat/completions"

	// 调用外部 API
	var response watermarkModel.ChatCompletionResponse
	status, body, err := http.Post(dyuAPIURL, req, apiKey, &response)

	if err != nil {
		log.Printf("watermarkService.RemoveWatermark 调用外部API失败, error: %v", err)
		return nil, err
	}

	// 即使状态码不是200，也返回解析后的响应体（外部API可能返回错误信息，但格式一致）
	if status != 200 {
		log.Printf("watermarkService.RemoveWatermark 外部API返回错误状态码: %d, body: %s", status, body)
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
