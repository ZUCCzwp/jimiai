package payRepo

import (
	"jiyu/global"
	"jiyu/model/payModel"
	"time"

	"gorm.io/gorm"
)

func SavePaymentLog(log payModel.PaymentLog) error {
	return global.DB.Model(&payModel.PaymentLog{}).
		Create(&log).Error
}

func FindByOrderNo(orderNo string) (payModel.PaymentLog, error) {
	var log payModel.PaymentLog
	err := global.DB.Model(&payModel.PaymentLog{}).
		Where("order_no = ?", orderNo).
		First(&log).Error
	return log, err
}

func UpdatePaymentLog(tx *gorm.DB, log payModel.PaymentLog) error {
	return tx.Model(&payModel.PaymentLog{}).
		Select("order_status", "trade_no").
		Where("order_no = ?", log.OrderNo).
		Updates(&log).Error
}

func ListForAdmin(page, limit, uid, orderStatus int, startTime, endTime string) ([]payModel.PaymentLog, error) {
	var results []payModel.PaymentLog

	tb := global.DB.Model(&payModel.PaymentLog{})
	if uid > 0 {
		tb = tb.Where("uid = ?", uid)
	}

	if startTime != "" {
		tb = tb.Where("created_at >= ?", startTime)
	}

	if endTime != "" {
		tb = tb.Where("created_at <= ?", endTime)
	}

	if orderStatus >= 0 {
		tb = tb.Where("order_status = ?", orderStatus)
	}

	err := tb.Order("id desc").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&results).Error
	return results, err
}

func Count(uid, orderStatus int, startTime, endTime string) (int, error) {
	var count int64
	tb := global.DB.Model(&payModel.PaymentLog{})
	if uid > 0 {
		tb = tb.Where("uid = ?", uid)
	}

	if orderStatus >= 0 {
		tb = tb.Where("order_status = ?", orderStatus)
	}

	if startTime != "" {
		tb = tb.Where("created_at >= ?", startTime)
	}

	if endTime != "" {
		tb = tb.Where("created_at <= ?", endTime)
	}
	err := tb.Count(&count).Error
	return int(count), err
}

func FindCurrentOnePayedByUid(uid int) (payModel.PaymentLog, error) {
	var results payModel.PaymentLog
	err := global.DB.Model(payModel.PaymentLog{}).
		Where("uid = ? AND order_status = ?", uid, payModel.PaymentStatusSuccess).
		Order("id desc").
		First(&results).Error
	return results, err
}

func FindCurrentOnePayedByUids(uids []string) ([]payModel.PaymentLog, error) {
	var results []payModel.PaymentLog
	err := global.DB.Model(payModel.PaymentLog{}).
		Where("uid in ? AND order_status = ?", uids, payModel.PaymentStatusSuccess).
		Order("id desc").
		Find(&results).Error
	return results, err
}

// FindTodayRMB 获取今日收入
func FindTodayRMB() (float64, error) {
	var results []payModel.PaymentLog
	err := global.DB.Model(payModel.PaymentLog{}).
		Where("created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", 1).
		Find(&results).Error

	sum := 0.0
	for _, v := range results {
		sum += v.Rmb
	}

	return sum, err
}

func SaveWithdrawalLog(log payModel.WithdrawalLog) error {
	return global.DB.Model(&payModel.WithdrawalLog{}).Create(&log).Error
}

func WithdrawalListForAdmin(page, limit, uid, status int, startTime, endTime string) ([]payModel.WithdrawalLog, error) {
	var results []payModel.WithdrawalLog
	tb := global.DB.Model(&payModel.WithdrawalLog{})
	if uid > 0 {
		tb = tb.Where("uid = ?", uid)
	}

	if startTime != "" {
		tb = tb.Where("created_at >= ?", startTime)
	}

	if endTime != "" {
		tb = tb.Where("created_at <= ?", endTime)
	}

	if status >= 0 {
		tb = tb.Where("status = ?", status)
	}

	err := tb.Order("id desc").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&results).Error
	return results, err
}

func CountWithdrawal(uid, status int, startTime, endTime string) (int, error) {
	var count int64
	tb := global.DB.Model(payModel.WithdrawalLog{})
	if uid > 0 {
		tb = tb.Where("uid = ?", uid)
	}

	if status >= 0 {
		tb = tb.Where("status = ?", status)
	}

	if startTime != "" {
		tb = tb.Where("created_at >= ?", startTime)
	}

	if endTime != "" {
		tb = tb.Where("created_at <= ?", endTime)
	}

	err := tb.Count(&count).Error
	return int(count), err
}

func EditWithdrawalStatus(id, status int) error {
	m := make(map[string]any)
	m["status"] = status
	m["handle_time"] = time.Now()

	return global.DB.Model(payModel.WithdrawalLog{}).Where("id = ?", id).Updates(m).Error
}

func FindTodayWithdrawalCount(uid int) (int, error) {
	var count int64
	err := global.DB.Model(payModel.WithdrawalLog{}).
		Where("uid = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", uid, 1).
		Count(&count).Error
	return int(count), err
}

func FindListByUid(uid int) ([]payModel.WithdrawalLog, error) {
	var results []payModel.WithdrawalLog
	err := global.DB.Model(payModel.WithdrawalLog{}).
		Where("uid = ?", uid).
		Order("id desc").Find(&results).Error
	return results, err
}
