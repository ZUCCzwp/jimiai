package cron

import (
	"fmt"
	"jiyu/repo/serverRepo"
	"jiyu/service/userService"
	"log"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
)

var (
	hourEvent = make([]func(), 0)
	dayEvent  = make([]func(), 0)
)

func init() {
	AddDayEvent(updateHotUser)
	AddDayEvent(resetCounter)
}

func Start() {
	c := cron.New()

	_, _ = c.AddFunc("@hourly", func() {
		for i := range hourEvent {
			hourEvent[i]()
		}
	})

	_, _ = c.AddFunc("@daily", func() {
		for i := range hourEvent {
			dayEvent[i]()
		}
	})

	c.Start()
}

func AddHourEvent(f func()) {
	hourEvent = append(hourEvent, f)
}

func AddDayEvent(f func()) {
	dayEvent = append(dayEvent, f)
}

// 更新热门用户
func updateHotUser() {
	err := userService.SelectHotUsersTable()
	if err != nil {
		log.Println("更新热门用户失败: ", err)
	}
}

// 重置用户统计次数
func resetCounter() {
	// 重置用户统计次数
	serverConfig := serverRepo.FindConfig()

	if carbon.FromStdTime(serverConfig.LastClearTime).IsToday() {
		fmt.Println("检查重置时间, 跳过")
		return
	}

	fmt.Println("重置用户统计次数")

	serverConfig.LastClearTime = time.Now()
	err := serverRepo.SaveConfig(serverConfig)
	if err != nil {
		log.Println("更新LastClearTime失败: ", err)
	}
}
