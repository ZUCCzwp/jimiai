package taskRepo

import (
	"jiyu/global"
	"jiyu/model/taskModel"
	"time"
)

// CreateTask 创建任务
func CreateTask(task *taskModel.Task) error {
	return global.DB.Model(&taskModel.Task{}).Create(task).Error
}

// GetTaskByTaskId 根据任务ID获取任务
func GetTaskByTaskId(taskId string) (*taskModel.Task, error) {
	var task taskModel.Task
	err := global.DB.Model(&taskModel.Task{}).Where("task_id = ?", taskId).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetTaskByUserId 根据用户ID获取任务列表
func GetTaskByUserId(userId int) ([]taskModel.Task, error) {
	var tasks []taskModel.Task
	err := global.DB.Model(&taskModel.Task{}).Where("user_id = ?", userId).Order("id DESC").Find(&tasks).Error
	return tasks, err
}

// UpdateTask 更新任务
func UpdateTask(taskId string, updates map[string]interface{}) error {
	return global.DB.Model(&taskModel.Task{}).Where("task_id = ?", taskId).Updates(updates).Error
}

// GetTaskList 获取任务列表
func GetTaskList(userId int, page, pageSize int, search, status, startTime, endTime string) ([]taskModel.TaskListResponse, error) {
	var results []taskModel.TaskListResponse

	query := global.DB.Model(&taskModel.Task{}).
		Select("id, task_id, model_name as model, video_type, prompt, status, progress, submit_time, end_time, time_spent, video_path, remark").
		Where("user_id = ?", userId)

	// 搜索关键词筛选（任务ID或提示词）
	if search != "" {
		query = query.Where("(task_id LIKE ? OR prompt LIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 时间范围筛选
	if startTime != "" && endTime != "" {
		start, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err == nil {
			end, err := time.Parse("2006-01-02 15:04:05", endTime)
			if err == nil {
				query = query.Where("submit_time >= ? AND submit_time <= ?", start, end)
			}
		}
	}

	// 分页和排序
	err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&results).Error

	return results, err
}

// CountTaskList 统计任务列表总数
func CountTaskList(userId int, search, status, startTime, endTime string) (int, error) {
	var count int64

	query := global.DB.Model(&taskModel.Task{}).
		Where("user_id = ?", userId)

	// 搜索关键词筛选（任务ID或提示词）
	if search != "" {
		query = query.Where("(task_id LIKE ? OR prompt LIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 时间范围筛选
	if startTime != "" && endTime != "" {
		start, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err == nil {
			end, err := time.Parse("2006-01-02 15:04:05", endTime)
			if err == nil {
				query = query.Where("submit_time >= ? AND submit_time <= ?", start, end)
			}
		}
	}

	err := query.Count(&count).Error
	return int(count), err
}

