package payService

import (
	"errors"
	"fmt"
	"jiyu/model/contextModel"
	"jiyu/model/payModel"
	"jiyu/model/userModel"
	"jiyu/repo/payRepo"
	"jiyu/repo/settingRepo"
	"log"
	"time"
)

func ListForAdmin(ctx contextModel.AdminContext, page, limit, uid, orderStatus int, startTime, endTime string) ([]payModel.PaymentLog, int, error) {
	list, err := payRepo.ListForAdmin(page, limit, uid, orderStatus, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	count, err := payRepo.Count(uid, orderStatus, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func WithdrawalListForAdmin(ctx contextModel.AdminContext, page, limit, uid, status int, startTime, endTime string) ([]payModel.WithdrawalLog, int, error) {
	list, err := payRepo.WithdrawalListForAdmin(page, limit, uid, status, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	count, err := payRepo.CountWithdrawal(uid, status, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func EditWithdrawal(ctx contextModel.AdminContext, id, status int) error {
	return payRepo.EditWithdrawalStatus(id, status)
}

func AddAccount(ctx contextModel.Context, data userModel.PaymentInfoRequest) error {
	return errors.New("支付账号功能已移除")
}

func DeleteAccount(ctx contextModel.Context, _uuid string) error {
	return errors.New("支付账号功能已移除")
}

func Withdrawal(ctx contextModel.Context, data payModel.WithdrawalRequest) error {
	withdrawalCount, err := payRepo.FindTodayWithdrawalCount(int(ctx.User.ID))
	if err != nil {
		return err
	}

	if withdrawalCount >= 1 {
		return errors.New("每天最多提现1次")
	}

	setting, err := settingRepo.Find()
	if err != nil {
		return err
	}

	if data.RMB < setting.Withdraw.MinWithdrawAmount {
		return errors.New(fmt.Sprintf("最低提现金额为%v元", setting.Withdraw.MinWithdrawAmount))
	}

	needCoin := int(data.RMB * setting.Withdraw.WithdrawRate)

	if ctx.User.JimiCoin < needCoin {
		return errors.New("余额不足")
	}

	ctx.User.JimiCoin -= needCoin

	wLog := payModel.WithdrawalLog{
		Uid:               int(ctx.User.ID),
		Rmb:               data.RMB,
		Ticket:            0,
		WithdrawalType:    data.AccountUUID, // 临时使用 AccountUUID 作为支付类型
		WithdrawalName:    "",
		WithdrawalAccount: "",
		Status:            0,
		HandleTime:        time.Now(),
	}

	err = payRepo.SaveWithdrawalLog(wLog)
	if err != nil {
		log.Printf("提现失败: %v", err)
		return errors.New("提现失败")
	}

	return nil
}

func WithdrawalList(ctx contextModel.Context) ([]payModel.WithdrawalListResponse, error) {
	logs, err := payRepo.FindListByUid(int(ctx.User.ID))
	if err != nil {
		log.Printf("获取提现记录失败: %v", err)
		return nil, errors.New("获取提现记录失败")
	}

	results := make([]payModel.WithdrawalListResponse, 0)
	for _, v := range logs {

		status := "未知"
		switch v.Status {
		case 0:
			status = "处理中"
		case 1:
			status = "成功"
		case 2:
			status = "已拒绝"
		}

		results = append(results, payModel.WithdrawalListResponse{
			Id:      int(v.ID),
			RMB:     v.Rmb,
			Account: v.WithdrawalAccount,
			Time:    v.CreatedAt,
			Status:  status,
		})
	}

	return results, nil
}
