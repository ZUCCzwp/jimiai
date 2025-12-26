package watermarkService

import (
	"jiyu/config"
	"jiyu/model/watermarkModel"
	"jiyu/util/http"
	"log"
)

const (
	defaultDyuAPIURL = "https://api.dyuapi.com/v1/chat/completions"
)

// RemoveWatermark 调用去水印接口
func RemoveWatermark(req watermarkModel.ChatCompletionRequest, apiKey string) (*watermarkModel.ChatCompletionResponse, error) {
	// 验证 API Key
	if apiKey == "" {
		log.Println("watermarkService.RemoveWatermark API Key 未配置")
		return nil, &ServiceError{Message: "API Key 未配置"}
	}

	// 从配置文件获取 API URL，如果没有配置则使用默认值
	dyuAPIURL := config.AppConfig.DyuAPIURL
	if dyuAPIURL == "" {
		dyuAPIURL = defaultDyuAPIURL
	}

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
