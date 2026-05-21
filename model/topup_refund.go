package model

import "github.com/QuantumNous/new-api/common"

const (
	RefundStatusPending    = "pending"
	RefundStatusProcessing = "processing"
	RefundStatusSuccess    = "success"
	RefundStatusFailed     = "failed"
)

type TopUpRefund struct {
	Id              int     `json:"id" gorm:"primaryKey"`
	TopUpId         int     `json:"topup_id" gorm:"index"`
	UserId          int     `json:"user_id" gorm:"index"`
	TradeNo         string  `json:"trade_no" gorm:"type:varchar(255);index"`
	RefundNo        string  `json:"refund_no" gorm:"type:varchar(128);uniqueIndex"`
	ProviderOrderNo string  `json:"provider_order_no" gorm:"type:varchar(128)"`
	RefundAmount    float64 `json:"refund_amount"`
	TotalAmount     float64 `json:"total_amount"`
	Reason          string  `json:"reason" gorm:"type:varchar(500)"`
	Status          string  `json:"status" gorm:"type:varchar(20)"`
	QuotaDeducted   int     `json:"quota_deducted"`
	OperatorId      int     `json:"operator_id"`
	CreateTime      int64   `json:"create_time"`
	CompleteTime    int64   `json:"complete_time"`
}

func (r *TopUpRefund) Insert() error {
	return DB.Create(r).Error
}

func (r *TopUpRefund) Update() error {
	return DB.Save(r).Error
}

func GetTopUpRefundByRefundNo(refundNo string) *TopUpRefund {
	var refund TopUpRefund
	err := DB.Where("refund_no = ?", refundNo).First(&refund).Error
	if err != nil {
		return nil
	}
	return &refund
}

func GetTopUpRefundsByTopUpId(topUpId int) []*TopUpRefund {
	var refunds []*TopUpRefund
	DB.Where("topup_id = ?", topUpId).Order("id desc").Find(&refunds)
	return refunds
}

func GetAllOpenPaymentRefunds(page, pageSize int) ([]*TopUpRefund, int64) {
	var refunds []*TopUpRefund
	var total int64
	DB.Model(&TopUpRefund{}).Count(&total)
	DB.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&refunds)
	return refunds, total
}

func CompleteOpenPaymentRefund(refundNo string) error {
	refund := GetTopUpRefundByRefundNo(refundNo)
	if refund == nil {
		return nil
	}
	if refund.Status == RefundStatusSuccess {
		return nil
	}

	refund.Status = RefundStatusSuccess
	refund.CompleteTime = common.GetTimestamp()

	if err := refund.Update(); err != nil {
		return err
	}

	topUp := GetTopUpByTradeNo(refund.TradeNo)
	if topUp == nil {
		return nil
	}

	topUp.RefundAmount += refund.RefundAmount
	if topUp.RefundAmount >= topUp.Money {
		topUp.RefundStatus = "refunded"
	} else {
		topUp.RefundStatus = "refunding"
	}
	_ = topUp.Update()

	quotaToDeduct := refund.QuotaDeducted
	if quotaToDeduct > 0 {
		_ = DecreaseUserQuota(refund.UserId, quotaToDeduct, true)
	}

	return nil
}

func FailOpenPaymentRefund(refundNo string) error {
	refund := GetTopUpRefundByRefundNo(refundNo)
	if refund == nil {
		return nil
	}
	if refund.Status == RefundStatusSuccess || refund.Status == RefundStatusFailed {
		return nil
	}
	refund.Status = RefundStatusFailed
	refund.CompleteTime = common.GetTimestamp()
	return refund.Update()
}
