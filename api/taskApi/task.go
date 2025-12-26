package taskApi

import (
	"fmt"
	"jiyu/model/contextModel"
	"jiyu/model/taskModel"
	"jiyu/service/taskService"
	"jiyu/util/response"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTask 创建任务接口
func CreateTask(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	var req taskModel.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	// 验证必填参数
	if req.TaskId == "" {
		response.Error(c, response.INVALID_PARAMS, "任务ID不能为空", nil)
		return
	}

	err := taskService.CreateTask(ctx, req)
	if err != nil {
		log.Printf("taskApi.CreateTask 创建任务失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "创建成功", nil)
}

// UpdateTask 更新任务接口
func UpdateTask(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	// 从路径参数获取任务ID
	taskId := c.Param("taskId")
	if taskId == "" {
		response.Error(c, response.INVALID_PARAMS, "任务ID不能为空", nil)
		return
	}

	var req taskModel.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("更新任务信息err:", err)
		response.Error(c, response.INVALID_PARAMS, "参数错误", nil)
		return
	}

	err := taskService.UpdateTask(ctx, taskId, req)
	if err != nil {
		log.Printf("taskApi.UpdateTask 更新任务失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "更新成功", nil)
}

// GetTaskList 获取任务列表接口
func GetTaskList(c *gin.Context) {
	ctx := c.MustGet("context").(contextModel.Context)

	// 解析分页参数
	pn := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pn, 10, 64)
	if err != nil || page < 1 {
		response.Error(c, response.INVALID_PARAMS, "参数错误：page必须是正整数", nil)
		return
	}

	ps := c.DefaultQuery("pageSize", "100")
	pageSize, err := strconv.ParseInt(ps, 10, 64)
	if err != nil || pageSize < 1 {
		response.Error(c, response.INVALID_PARAMS, "参数错误：pageSize必须是正整数", nil)
		return
	}

	// 解析筛选参数
	search := c.Query("search")
	status := c.Query("status")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	req := taskModel.TaskListRequest{
		Page:      int(page),
		PageSize:  int(pageSize),
		Search:    search,
		Status:    status,
		StartTime: startTime,
		EndTime:   endTime,
	}

	// 调用service层获取数据
	list, count, err := taskService.GetTaskList(ctx, req)
	if err != nil {
		log.Printf("taskApi.GetTaskList 获取任务列表失败, error: %v", err)
		response.Error(c, response.ERROR, err.Error(), nil)
		return
	}

	response.Success(c, "ok", gin.H{
		"list":  list,
		"total": count,
	})
}
