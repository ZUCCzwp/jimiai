package videoModel

// VideoGenerationRequest 图生视频请求
type VideoGenerationRequest struct {
	Prompt  string `form:"prompt" binding:"required"`  // 提示词
	Model   string `form:"model" binding:"required"`   // 模型名称
	Size    string `form:"size" binding:"required"`    // 尺寸，如 "1280x720"
	Seconds string `form:"seconds" binding:"required"` // 秒数，如 "4"
}

// VideoGenerationResponse 图生视频响应（与外部API保持一致）
type VideoGenerationResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Model     string `json:"model"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	CreatedAt int64  `json:"created_at"`
	Seconds   string `json:"seconds"`
	Size      string `json:"size"`
}

// VideoQueryResponse 查询视频任务响应（与外部API保持一致）
type VideoQueryResponse struct {
	ID          string `json:"id"`
	Size        string `json:"size"`
	Model       string `json:"model"`
	Object      string `json:"object"`
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	VideoURL    string `json:"video_url"`
	CreatedAt   int64  `json:"created_at"`
	CompletedAt int64  `json:"completed_at"`
}
