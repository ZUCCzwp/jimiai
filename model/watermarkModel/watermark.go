package watermarkModel

// ChatCompletionRequest 聊天完成请求（与外部API保持一致）
type ChatCompletionRequest struct {
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// ChatCompletionResponse 聊天完成响应（与外部API保持一致）
type ChatCompletionResponse struct {
	ID             string   `json:"id"`
	Object         string   `json:"object"`
	Created        int64    `json:"created"`
	Model          string   `json:"model"`
	Choices        []Choice `json:"choices"`
	Usage          Usage    `json:"usage"`
	Links          *Links   `json:"links,omitempty"`
	PostID         string   `json:"post_id,omitempty"`
	OriginalInput  string   `json:"original_input,omitempty"`
	PostInfo       *PostInfo `json:"post_info,omitempty"`
}

// Choice 选择项
type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

// Usage 使用情况
type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
}

// Links 链接信息
type Links struct {
	Gif       string `json:"gif,omitempty"`
	Text      string `json:"text,omitempty"`
	ID        string `json:"id,omitempty"`
	Mp4       string `json:"mp4,omitempty"`
	Mp4Wm     string `json:"mp4_wm,omitempty"`
	Md        string `json:"md,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// PostInfo 帖子信息
type PostInfo struct {
	AttachmentsCount int    `json:"attachments_count"`
	ViewCount        int    `json:"view_count"`
	LikeCount        int    `json:"like_count"`
	Title            string `json:"title"`
}

