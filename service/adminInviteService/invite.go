package adminInviteService

import (
	"jiyu/model/inviteModel"
	"jiyu/model/payModel"
	"jiyu/model/userModel"
	"jiyu/repo/inviteRepo"
	"jiyu/repo/payRepo"
	"jiyu/repo/userRepo"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func AdminInviteList(page int, pageSize int) ([]inviteModel.InviteLogResponse, int64, error) {
	result, err := inviteRepo.AdminList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := inviteRepo.AdminCount()
	if err != nil {
		return nil, 0, err
	}

	var allUserIds []string
	for _, item := range result {
		allUserIds = append(allUserIds,
			cast.ToString(item.UserID),
			cast.ToString(item.InviteUserID))
	}
	allUserIds = lo.Uniq(allUserIds)

	users, err := userRepo.FindUserByIds(allUserIds)
	if err != nil {
		return nil, 0, err
	}

	userMap := make(map[int64]userModel.User)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}
	// 查询用户首次充值金额记录
	recharges, _ := payRepo.FindCurrentOnePayedByUids(allUserIds)

	userRechargeMap := make(map[int64]payModel.PaymentLog)
	for _, recharge := range recharges {
		userRechargeMap[int64(recharge.Uid)] = recharge
	}

	var response []inviteModel.InviteLogResponse
	for _, item := range result {
		if _, ok := userMap[item.InviteUserID]; !ok {
			continue
		}

		var fcm float64 = 0
		if _, ok := userRechargeMap[item.UserID]; ok {
			fcm = userRechargeMap[item.UserID].Rmb
		}
		response = append(response, inviteModel.InviteLogResponse{
			Id:             int64(item.ID),
			UserID:         item.UserID,
			InviteUserID:   item.InviteUserID,
			InviteType:     item.InviteType,
			InviteCoins:    item.InviteCoins,
			InviteTime:     userMap[item.InviteUserID].CreatedAt.Format(time.DateTime),
			FirstChargeNum: fcm,
		})
	}

	return response, count, nil
}

func InviteSuperior(uid int, page int, limit int) ([]inviteModel.InviteLogResponse, int64, error) {
	result, err := inviteRepo.AdminInviteSuperior(uid, page, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := inviteRepo.AdminInviteSuperiorCount(uid)
	if err != nil {
		return nil, 0, err
	}

	var allUserIds []string
	for _, item := range result {
		allUserIds = append(allUserIds,
			cast.ToString(item.UserID),
			cast.ToString(item.InviteUserID))
	}
	allUserIds = lo.Uniq(allUserIds)

	users, err := userRepo.FindUserByIds(allUserIds)
	if err != nil {
		return nil, 0, err
	}

	userMap := make(map[int64]userModel.User)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

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

	return response, count, nil
}

func InviteInferior(uid int, page int, limit int) ([]inviteModel.InviteLogResponse, int64, error) {
	result, err := inviteRepo.AdminInviteInferior(uid, page, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := inviteRepo.AdminInviteInferiorCount(uid)
	if err != nil {
		return nil, 0, err
	}

	var allUserIds []string
	for _, item := range result {
		allUserIds = append(allUserIds,
			cast.ToString(item.UserID),
			cast.ToString(item.InviteUserID))
	}
	allUserIds = lo.Uniq(allUserIds)

	users, err := userRepo.FindUserByIds(allUserIds)
	if err != nil {
		return nil, 0, err
	}

	userMap := make(map[int64]userModel.User)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

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

	return response, count, nil
}
