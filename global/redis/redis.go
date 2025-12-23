package redis

import (
	"context"
	"jiyu/config"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	KeyMomentDetail       = "jimiai:moment:detail"
	KeyMomentLiked        = "jimiai:moment:liked"
	KeyCommentLiked       = "jimiai:comment:liked"
	KeyMomentList         = "jimiai:moment:list"
	KeyBottleMomentDetail = "jimiai:bottle:moment:detail"
	KeyBottleLiked        = "jimiai:bottle:moment:liked"
	KeyBottleCommentLiked = "jimiai:bottle:comment:liked"
	KeySMSCode            = "jimiai:sms:code"
	KeyEmailCode          = "jimiai:email:code"
	KeyBottleRead         = "jimiai:bottle:read"
	KeyHotUsers           = "jimiai:user:hot"
	KeyUserRanking        = "jimiai:invite:user:ranking" //邀请排行
	KeyTokenBlacklist     = "jimiai:token:blacklist"     //token黑名单
)

var (
	RDB *redis.Client
	Ctx context.Context
)

func NewRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Addr,
		Password: config.RedisConfig.Pwd,
		DB:       config.RedisConfig.DB,
		//连接池容量及闲置连接数量
		PoolSize:     15,              // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,              //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。

	})

	Ctx = context.Background()

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Printf("redis.NewRedis ping error: %v\n", err)
		os.Exit(-1)
	}
}
