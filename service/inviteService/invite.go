package inviteService

import (
	"jiyu/global"
	"jiyu/global/redis"
	"jiyu/model/inviteModel"
	"jiyu/model/userModel"
	"jiyu/repo/inviteRepo"
	"jiyu/repo/userRepo"
	"log"
	"strconv"
	"time"

	redis2 "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// GetUserRanking 获取用户排行榜
func GetUserRanking() ([]*inviteModel.UserRank, error) {
	data, err := getLeaderboard(-1)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetUserRankingByPage(page int, pageSize int) ([]*inviteModel.UserRank, int64, error) {
	// 计算起始位置和结束位置
	start := (page - 1) * pageSize
	stop := start + pageSize - 1

	// 获取排行榜分页数据
	data, err := redis.RDB.ZRevRangeWithScores(redis.Ctx, redis.KeyUserRanking, int64(start), int64(stop)).Result()
	if err != nil {
		return nil, 0, err
	}

	// 获取用户ID列表
	var userIds []string
	for _, z := range data {
		userIds = append(userIds, z.Member.(string))
	}

	// 获取用户信息
	users, err := userRepo.FindUserByIds(userIds)
	if err != nil {
		return nil, 0, err
	}

	// 构建用户map
	userMap := make(map[string]userModel.User)
	for _, user := range users {
		userMap[strconv.Itoa(int(user.ID))] = user
	}

	// 构建排行榜数据
	var result []*inviteModel.UserRank
	for i, z := range data {
		userId := z.Member.(string)
		if user, ok := userMap[userId]; ok {
			result = append(result, &inviteModel.UserRank{
				Rank:     int64(start + i + 1),
				UserId:   int64(user.ID),
				Score:    z.Score,
				Coins:    int64(user.JimiCoin),
				JoinTime: user.CreatedAt.Unix(),
			})
		}
	}

	count, err := redis.RDB.ZCard(redis.Ctx, redis.KeyUserRanking).Result()
	if err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

// GetTopNLeaderboard 获取前N名排行榜
func GetTopNLeaderboard(n int) ([]*inviteModel.UserRank, error) {
	return getLeaderboard(n)
}

func GetUserRankByUserId(userId int64) (*inviteModel.UserRank, error) {
	//获取用户排名
	rank, err := getUserRankingByUserId(userId)
	if err != nil {
		return nil, err
	}
	//获取用户score分数
	score, err := redis.RDB.ZScore(redis.Ctx, redis.KeyUserRanking, strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return nil, err
	}

	//获取用户信息
	user, err := userRepo.FindByID(int(userId))
	if err != nil {
		return nil, err
	}

	return &inviteModel.UserRank{
		Rank:     rank,
		Coins:    int64(user.JimiCoin),
		JoinTime: user.CreatedAt.Unix(),
		UserId:   int64(user.ID),
		Score:    score,
	}, nil
}

// 获取某个用户的排名
func getUserRankingByUserId(userId int64) (int64, error) {
	result, err := redis.RDB.ZRevRank(redis.Ctx, redis.KeyUserRanking, strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return 0, err
	}

	// 排行榜从1开始
	return result + 1, nil
}

// 获取排行榜通用方法
func getLeaderboard(n int) ([]*inviteModel.UserRank, error) {
	var (
		err      error
		rankings []redis2.Z
	)

	if n < 0 {
		// 获取全部排行
		rankings, err = redis.RDB.ZRevRangeWithScores(redis.Ctx, redis.KeyUserRanking, 0, -1).Result()
	} else {
		// 获取前N名
		rankings, err = redis.RDB.ZRangeWithScores(redis.Ctx, redis.KeyUserRanking, 0, int64(n-1)).Result()
	}

	if err != nil {
		log.Printf("获取用户排行榜失败: %v", err)
		return nil, err
	}

	ranks := make([]*inviteModel.UserRank, 0, len(rankings))
	users, err := userRepo.ListAll()
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		return nil, err
	}
	userMap := lo.SliceToMap(users, func(user userModel.User) (int64, userModel.User) {
		return int64(user.ID), user
	})
	for i, ranking := range rankings {
		userId := ranking.Member.(string)
		score := ranking.Score

		if _, ok := userMap[cast.ToInt64(userId)]; !ok {
			continue
		}

		user := userMap[cast.ToInt64(userId)]

		// 如果用户没有机遇币，则跳过
		if user.JimiCoin == 0 {
			continue
		}
		ranks = append(ranks, &inviteModel.UserRank{
			UserId:   cast.ToInt64(userId),
			Score:    score,
			Coins:    int64(user.JimiCoin),
			JoinTime: user.CreatedAt.Unix(),
			Rank:     int64(i + 1),
		})
	}

	return ranks, nil
}

// UpdateUserScore 计算用户得分并更新到Redis
// InviteRewardValue根据配置表获取更新
func UpdateUserScore(userId int64, coins int) error {
	user, err := userRepo.FindByID(int(userId))
	if err != nil {
		log.Printf("获取用户详细信息失败: %v", err)
		return err
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户的邀请人数
	if err = userRepo.UpdateInviteCountTx(tx, user); err != nil {
		tx.Rollback()
		log.Printf("更新用户邀请人数失败: %v", err)
		return err
	}

	user.JimiCoin = coins
	if err = userRepo.UpdateJimiCoinTx(tx, user); err != nil {
		tx.Rollback()
		log.Printf("更新用户信息失败: %v", err)
		return err
	}

	// Redis更新
	sumCoins := coins + user.JimiCoin
	score := calculateScore(sumCoins, user.CreatedAt)
	pipe := redis.RDB.Pipeline()
	pipe.ZAdd(redis.Ctx, redis.KeyUserRanking, redis2.Z{
		Score:  score,
		Member: userId,
	})

	// 执行Redis pipeline和MySQL事务
	_, err = pipe.Exec(redis.Ctx)
	if err != nil {
		tx.Rollback()
		log.Printf("更新Redis失败: %v", err)
		return err
	}

	// 提交MySQL事务
	if err = tx.Commit().Error; err != nil {
		// 如果MySQL事务提交失败,需要回滚Redis
		redis.RDB.ZRem(redis.Ctx, redis.KeyUserRanking, userId)
		log.Printf("提交事务失败: %v", err)
		return err
	}

	return nil
}

// 计算用户分数 (机遇币 + 注册时间)
func calculateScore(opportunityCoin int, registrationTime time.Time) float64 {
	return float64(opportunityCoin) + float64(time.Now().Unix()-registrationTime.Unix())/1000000
}

func GetInviteRecord(userId int64, page int, pageSize int) ([]inviteModel.InviteLogResponse, error) {
	result, err := inviteRepo.ListByUserID(userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	inviteUserIds := lo.Uniq(lo.Map(result, func(item inviteModel.InviteLog, index int) string {
		return cast.ToString(item.InviteUserID)
	}))

	users, err := userRepo.FindUserByIds(inviteUserIds)
	if err != nil {
		return nil, err
	}
	userMap := lo.SliceToMap(users, func(user userModel.User) (int64, userModel.User) {
		return int64(user.ID), user
	})

	var response []inviteModel.InviteLogResponse
	for _, item := range result {
		if _, ok := userMap[item.InviteUserID]; !ok {
			continue
		}
		response = append(response, inviteModel.InviteLogResponse{
			Id:           int64(item.ID),
			UserID:       item.UserID,
			InviteUserID: item.InviteUserID,
			InviteType:   item.InviteType,
			InviteCoins:  item.InviteCoins,
			InviteTime:   userMap[item.InviteUserID].CreatedAt.Format(time.DateTime),
		})
	}

	return response, nil
}
