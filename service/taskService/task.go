package taskService

import (
	"errors"
	"jiyu/model/contextModel"
	"jiyu/model/taskModel"
	"jiyu/repo/taskRepo"
	"log"
	"time"
)

// CreateTask 创建任务
func CreateTask(ctx contextModel.Context, req taskModel.CreateTaskRequest) error {
	userId := int(ctx.User.ID)

	// 验证必填参数
	if req.TaskId == "" {
		return errors.New("任务ID不能为空")
	}

	// 检查任务ID是否已存在
	existingTask, err := taskRepo.GetTaskByTaskId(req.TaskId)
	if err == nil && existingTask != nil {
		return errors.New("任务ID已存在")
	}

	// 解析提交时间
	var submitTime time.Time
	if req.SubmitTime != "" {
		submitTime, err = time.Parse("2006-01-02 15:04:05", req.SubmitTime)
		if err != nil {
			// 尝试其他时间格式
			submitTime, err = time.Parse(time.RFC3339, req.SubmitTime)
			if err != nil {
				log.Printf("taskService.CreateTask 解析提交时间失败, 使用当前时间, error: %v", err)
				submitTime = time.Now()
			}
		}
	} else {
		submitTime = time.Now()
	}

	// 设置默认状态
	status := req.Status
	if status == "" {
		status = "pending"
	}

	// 验证进度范围
	progress := req.Progress
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	task := &taskModel.Task{
		UserId:     userId,
		TaskId:     req.TaskId,
		ModelName:  req.Model,
		VideoType:  req.VideoType,
		Prompt:     req.Prompt,
		Status:     status,
		Progress:   progress,
		SubmitTime: submitTime,
	}

	err = taskRepo.CreateTask(task)
	if err != nil {
		log.Printf("taskService.CreateTask 创建任务失败, error: %v", err)
		return errors.New("创建任务失败")
	}

	return nil
}

// UpdateTask 更新任务
func UpdateTask(ctx contextModel.Context, taskId string, req taskModel.UpdateTaskRequest) error {
	userId := int(ctx.User.ID)

	// 验证任务ID
	if taskId == "" {
		return errors.New("任务ID不能为空")
	}

	// 检查任务是否存在且属于当前用户
	task, err := taskRepo.GetTaskByTaskId(taskId)
	if err != nil {
		return errors.New("任务不存在")
	}

	// 验证任务是否属于当前用户
	if task.UserId != userId {
		return errors.New("无权操作此任务")
	}

	// 构建更新数据
	updates := make(map[string]interface{})

	// 更新状态
	if req.Status != "" {
		updates["status"] = req.Status
	}

	// 更新进度
	if req.Progress >= 0 && req.Progress <= 100 {
		updates["progress"] = req.Progress
	} else if req.Progress > 100 {
		updates["progress"] = 100
	} else if req.Progress < 0 {
		updates["progress"] = 0
	}

	// 更新结束时间
	if req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err != nil {
			// 尝试其他时间格式
			endTime, err = time.Parse(time.RFC3339, req.EndTime)
			if err != nil {
				log.Printf("taskService.UpdateTask 解析结束时间失败, error: %v", err)
			} else {
				updates["end_time"] = endTime
			}
		} else {
			updates["end_time"] = endTime
		}
	}

	// 更新花费时间（秒数）
	if req.TimeSpent > 0 {
		updates["time_spent"] = req.TimeSpent
	}

	// 更新视频路径
	if req.VideoPath != "" {
		updates["video_path"] = req.VideoPath
	}

	// 更新备注
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	// 执行更新
	err = taskRepo.UpdateTask(taskId, updates)
	if err != nil {
		log.Printf("taskService.UpdateTask 更新任务失败, error: %v", err)
		return errors.New("更新任务失败")
	}

	return nil
}

// GetTaskList 获取任务列表
func GetTaskList(ctx contextModel.Context, req taskModel.TaskListRequest) ([]taskModel.TaskListResponse, int, error) {
	userId := int(ctx.User.ID)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 100
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	list, err := taskRepo.GetTaskList(userId, req.Page, req.PageSize, req.Search, req.Status, req.StartTime, req.EndTime)
	if err != nil {
		log.Printf("taskService.GetTaskList 获取任务列表失败, error: %v", err)
		return nil, 0, errors.New("获取任务列表失败")
	}

	count, err := taskRepo.CountTaskList(userId, req.Search, req.Status, req.StartTime, req.EndTime)
	if err != nil {
		log.Printf("taskService.GetTaskList 统计任务列表总数失败, error: %v", err)
		return nil, 0, errors.New("统计任务列表总数失败")
	}

	return list, count, nil
}

