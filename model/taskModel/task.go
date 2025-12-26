package taskModel

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务表
type Task struct {
	gorm.Model
	UserId     int        `json:"user_id" gorm:"type:int;not null;default:0;comment:'用户ID';index:idx_user_id"`                 // 用户ID
	TaskId     string     `json:"task_id" gorm:"type:varchar(100);not null;default:'';comment:'任务ID';uniqueIndex:idx_task_id"` // 任务ID
	ModelName  string     `json:"model" gorm:"type:varchar(100);default:'';comment:'模型名称';index:idx_model"`                    // 模型名称
	VideoType  string     `json:"video_type" gorm:"type:varchar(50);default:'';comment:'视频类型（图生视频）'"`                          // 视频类型
	Prompt     string     `json:"prompt" gorm:"type:varchar(1000);default:'';comment:'提示词'"`                                   // 提示词
	Status     string     `json:"status" gorm:"type:varchar(50);not null;default:'pending';comment:'任务状态';index:idx_status"`   // 任务状态
	Progress   int        `json:"progress" gorm:"type:int;not null;default:0;comment:'进度(0-100)'"`                             // 进度
	SubmitTime time.Time  `json:"submit_time" gorm:"type:datetime;comment:'提交时间'"`                                             // 提交时间
	EndTime    *time.Time `json:"end_time" gorm:"type:datetime;comment:'结束时间'"`                                                // 结束时间
	TimeSpent  int        `json:"time_spent" gorm:"type:int;default:0;comment:'花费时间(秒数)'"`                                     // 花费时间(秒数)
	VideoPath  string     `json:"video_path" gorm:"type:varchar(500);default:'';comment:'视频路径'"`                               // 视频路径
	Remark     string     `json:"remark" gorm:"type:varchar(500);default:'';comment:'备注'"`                                     // 备注
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskId     string `json:"taskId" binding:"required"` // 任务ID
	Model      string `json:"model"`                     // 模型名称
	VideoType  string `json:"videoType"`                 // 视频类型（图生视频）
	Prompt     string `json:"prompt"`                    // 提示词
	Status     string `json:"status"`                    // 任务状态
	Progress   int    `json:"progress"`                  // 进度
	SubmitTime string `json:"submitTime"`                // 提交时间
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Status    string `json:"status"`    // 任务状态
	Progress  int    `json:"progress"`  // 进度
	EndTime   string `json:"endTime"`   // 结束时间
	TimeSpent int    `json:"timeSpent"` // 花费时间秒数
	VideoPath string `json:"videoPath"` // 视频路径
	Remark    string `json:"remark"`    // 备注
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	Search    string `json:"search"`     // 搜索关键词（任务ID或提示词）
	Status    string `json:"status"`     // 任务状态
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	ID         uint       `json:"id"`         // 序号
	TaskId     string     `json:"taskId"`     // 任务ID
	Model      string     `json:"model"`      // 模型名称
	VideoType  string     `json:"videoType"`  // 视频类型
	Prompt     string     `json:"prompt"`     // 提示词
	Status     string     `json:"status"`     // 任务状态
	Progress   int        `json:"progress"`   // 进度
	SubmitTime time.Time  `json:"submitTime"` // 提交时间
	EndTime    *time.Time `json:"endTime"`    // 结束时间
	TimeSpent  int        `json:"timeSpent"`  // 花费时间(秒数)
	VideoPath  string     `json:"videoPath"`  // 视频路径
	Remark     string     `json:"remark"`     // 备注
}
