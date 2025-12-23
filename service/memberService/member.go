package memberService

import (
	"errors"
	"jiyu/model/contextModel"
	"jiyu/model/memberModel"
	"jiyu/model/userModel"
	"jiyu/repo/memberRepo"
	"jiyu/repo/payRepo"
	"jiyu/repo/userRepo"
	"strconv"
	"time"
)

func AdminMemberConfigList(ctx contextModel.AdminContext, page int, limit int) ([]memberModel.MemberConfig, int64, error) {
	count, err := memberRepo.CountMemberConfig()
	if err != nil {
		return nil, 0, err
	}

	memberList, err := memberRepo.ListPageMemberConfig(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return memberList, count, nil
}

func UpdateMemberConfig(ctx contextModel.AdminContext, member memberModel.MemberConfig) error {
	return memberRepo.UpdateMemberConfig(member)
}

func ListMemberConfig() ([]memberModel.MemberConfig, error) {
	return memberRepo.ListMemberConfig()
}

func AdminMemberList(ctx contextModel.AdminContext, page, limit int) ([]memberModel.MemberListResponse, int64, error) {
	count, err := userRepo.CountMember()
	if err != nil {
		return nil, 0, err
	}

	memberList, err := userRepo.ListMemberByPage(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// 根据返回的用户id查询用户充值记录返回最新购买时间
	rechargeMap := make(map[int64]time.Time)
	for _, v := range memberList {
		recharge, err := payRepo.FindCurrentOnePayedByUid(int(v.ID))
		if err != nil {
			continue
		}
		// 根据uid建立映射,key为uid，value为最新购买时间
		rechargeMap[int64(v.ID)] = recharge.CreatedAt
	}

	var memberListResponse []memberModel.MemberListResponse
	for _, member := range memberList {
		buyTime := ""
		if recharge, ok := rechargeMap[int64(member.ID)]; ok {
			buyTime = recharge.Format("2006-01-02 15:04:05")
		}
		vipTimeStr := ""
		if member.VipTime != nil {
			vipTimeStr = member.VipTime.Format("2006-01-02 15:04:05")
		}
		memberListResponse = append(memberListResponse, memberModel.MemberListResponse{
			ID:      int(member.ID),
			Phone:   member.Phone,
			BuyTime: buyTime,
			VipTime: vipTimeStr,
		})
	}

	return memberListResponse, count, nil
}

func EditMember(ctx contextModel.AdminContext, member memberModel.MemberEdit) error {
	// 查询用户是否存在
	uid, err := strconv.Atoi(member.Uid)
	if err != nil {
		return errors.New("用户ID格式错误")
	}

	m, err := userRepo.FindByID(uid)
	if err != nil {
		return errors.New("该用户不存在。")
	}

	vipTime, err := time.Parse("2006-01-02 15:04:05", member.VipTime)
	if err != nil {
		return err
	}

	// 比较会员到期时间 是否大于已存在的时间
	if m.VipTime != nil && m.VipTime.After(vipTime) {
		return errors.New("请仔细核对会员到期时间。")
	}

	user := userModel.User{
		VipTime: &vipTime,
	}
	return userRepo.EditMemberVipTime(uid, &user)
}
