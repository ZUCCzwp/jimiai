package main

import (
	"io"
	"jiyu/config"
	"jiyu/cron"
	"jiyu/global"
	"jiyu/global/redis"
	"jiyu/model/serverModel"
	"jiyu/repo/serverRepo"
	"jiyu/router"
	"jiyu/util/pay/wechat"
	"log"
	"math/rand"
	"os"
	"time"
)

// go build -o ./bin/app ./cmd/app/main.go
// nohup ./app > out.log 2>1&1 &
func main() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())

	// 初始化配置
	config.InitConfig()

	f := LogInit()
	defer f.Close()

	// 连接数据库
	_ = global.NewDB()

	// 初始化redis
	redis.NewRedis()

	// 自动迁移
	global.AutoMigrate()

	// 覆盖配置类
	// setting, err := settingRepo.Find()
	// if err != nil {
	// 	log.Fatalln("从数据库获取配置失败:", err)
	// }
	// setting.InitToConfig()

	// 初始化微信支付客户端
	global.WxClient = wechat.NewWechatClient()

	// 初始化服务器配置
	serverConfig := serverRepo.FindConfig()
	if serverConfig == nil || serverConfig.ID == 0 {
		s := serverModel.NewConfig()
		serverRepo.CreateConfig(s)
	}

	if config.LocalConfig.IsMaster {
		// 定时任务
		cron.Start()
	}

	r := router.NewRouter()
	_ = r.Run(config.LocalConfig.RunPort)
}

func LogInit() *os.File {
	// 打开日志文件
	f, err := os.OpenFile("./runtime/log/app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalln("日志文件打开失败")
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(f, os.Stdout)
	log.SetOutput(multiWriter)

	return f
}
