package main

import (
	"io"
	"jiyu/config"
	"jiyu/global"
	"jiyu/global/redis"
	"jiyu/model/adminUserModel"
	"jiyu/model/globalModel"
	"jiyu/model/settingModel"
	"jiyu/repo/adminUserRepo"
	"jiyu/repo/globalRepo"
	"jiyu/repo/settingRepo"
	"jiyu/router"
	"log"
	"math/rand"
	"os"
	"time"
)

// go build -o ../../bin/
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

	// 创建默认用户
	createDefaultUser()

	// 创建默认设置
	createDefaultSetting()

	// 创建默认 用户协议/隐私政策
	createDefaultDoc()

	// 覆盖配置类
	setting, err := settingRepo.Find()
	if err != nil {
		log.Fatalln("从数据库获取配置失败:", err)
	}
	setting.InitToConfig()

	r := router.NewAdminRouter()
	_ = r.Run(config.AdminAppConfig.RunPort)
}

func LogInit() *os.File {
	// 打开日志文件
	f, err := os.OpenFile("./runtime/log/admin.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalln("日志文件打开失败")
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(f, os.Stdout)
	log.SetOutput(multiWriter)

	return f
}

func createDefaultUser() {
	count, err := adminUserRepo.Count()
	if err == nil && count == 0 {
		err := adminUserRepo.Create(adminUserModel.NewUser())
		if err != nil {
			log.Println("创建默认用户失败:", err)
		} else {
			log.Println("创建默认用户成功")
		}
	}
}

func createDefaultSetting() {
	count, err := settingRepo.Count()
	if err == nil && count == 0 {
		err := settingRepo.Save(&settingModel.Setting{})
		if err != nil {
			log.Println("创建默认设置失败:", err)
		} else {
			log.Println("创建默认设置成功")
		}
	}
}

func createDefaultDoc() {
	count, err := globalRepo.DocCount()
	if err == nil && count == 0 {
		doc := globalModel.Doc{UserAgreement: "", PrivacyPolicy: ""}
		err := globalRepo.SaveDoc(&doc)
		if err != nil {
			log.Println("创建默认用户协议/隐私政策失败:", err)
		} else {
			log.Println("创建默认用户协议/隐私政策成功")
		}
	}
}
