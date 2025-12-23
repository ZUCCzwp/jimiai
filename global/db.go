package global

import (
	"context"
	"fmt"
	"jiyu/config"
	"jiyu/model/adminRouterModel"
	"jiyu/model/adminUserModel"
	"jiyu/model/appVersionModel"
	"jiyu/model/chargeModel"
	"jiyu/model/globalModel"
	"jiyu/model/inviteModel"
	"jiyu/model/memberModel"
	"jiyu/model/payModel"
	"jiyu/model/redeemCodeModel"
	"jiyu/model/serverModel"
	"jiyu/model/userModel"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func GetDB(ctx context.Context) *gorm.DB {
	traceId := ctx.Value(TraceIDKey{}).(string)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n"+traceId+" > ", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)

	return DB.Session(&gorm.Session{Logger: newLogger})
}

func NewDB() *gorm.DB {
	// 连接数据库. 如果数据库不存在，则创建数据库
	db, err := newDB(createDSN(config.DBConfig.DBName))

	if err != nil {
		log.Fatalf("数据库连接失败: %s\n", err.Error())
	}

	sqlDB, _ := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// 设置数据库闲链接超时时间
	sqlDB.SetConnMaxLifetime(300 * time.Second)

	DB = db
	return db
}

func AutoMigrate() {
	err := DB.AutoMigrate(&userModel.User{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&serverModel.Config{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&globalModel.Image{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&adminUserModel.User{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&payModel.PaymentLog{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&payModel.WithdrawalLog{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&globalModel.Feedback{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&globalModel.QA{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&globalModel.SMSCode{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.Set("gorm:table_options", " DEFAULT CHARSET=utf8mb4").AutoMigrate(&globalModel.Doc{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&globalModel.RefIdIp{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&adminRouterModel.Router{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %s\n", err.Error())
	}

	err = DB.AutoMigrate(&inviteModel.InviteLog{})
	if err != nil {
		log.Fatalf("InviteLog数据库迁移失败: %s\n", err.Error())
	}
	err = DB.AutoMigrate(&appVersionModel.AppVersion{})
	if err != nil {
		log.Fatalf("AppVersion数据库迁移失败: %s\n", err.Error())
	}
	err = DB.AutoMigrate(&chargeModel.ChargeInfo{})
	if err != nil {
		log.Fatalf("充值配置数据库迁移失败: %s\n", err.Error())
	}
	err = DB.AutoMigrate(&memberModel.MemberConfig{})
	if err != nil {
		log.Fatalf("会员配置数据库迁移失败: %s\n", err.Error())
	}
	err = DB.AutoMigrate(&redeemCodeModel.RedeemCode{})
	if err != nil {
		log.Fatalf("兑换码数据库迁移失败: %s\n", err.Error())
	}
}

var retryCount = 0

func newDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.DBConfig.TablePrefix,
		},
	})

	if err != nil && strings.Contains(err.Error(), "Unknown database") && retryCount < 3 {
		retryCount++
		log.Println("数据库不存在，尝试创建数据库")
		tempDB, err := newDB(createDSN("mysql"))

		if err != nil {
			return nil, err
		}

		err = tempDB.Exec("CREATE DATABASE " + config.DBConfig.DBName).Error

		if err != nil {
			return nil, err
		}

		log.Printf("数据库 %s 创建成功\n", config.DBConfig.DBName)

		return newDB(createDSN(config.DBConfig.DBName))
	}

	return db, err
}

func createDSN(dbName string) string {
	c := config.DBConfig
	timeZone := url.QueryEscape(c.TimeZone)

	if dbName == "" {
		dbName = c.DBName
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", c.UserName, c.Password, c.DbHost, c.DbPort, dbName, timeZone)

	return dsn
}
