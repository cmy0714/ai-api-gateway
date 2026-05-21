package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/service/paymentopen"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type OpenPaymentRequest struct {
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

func RequestOpenPayment(c *gin.Context) {
	if !isOpenPaymentTopUpEnabled() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "支付中心未启用"})
		return
	}

	var req OpenPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.Amount < setting.OpenPaymentMinTopUp {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %.2f", setting.OpenPaymentMinTopUp)})
		return
	}

	method := setting.FindOpenPaymentMethod(req.PaymentMethod, operation_setting.PayMethods)
	if method == nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "支付方式不存在"})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}
	payMoney := getPayMoney(req.Amount, group)
	if payMoney < 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	tradeNo := fmt.Sprintf("USR%dNO%s%d", id, common.GetRandomString(6), time.Now().Unix())

	amount := int64(req.Amount)
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dAmount := decimal.NewFromFloat(req.Amount)
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		amount = dAmount.Div(dQuotaPerUnit).IntPart()
	}

	topUp := &model.TopUp{
		UserId:          id,
		Amount:          amount,
		Money:           payMoney,
		TradeNo:         tradeNo,
		PaymentMethod:   req.PaymentMethod,
		PaymentProvider: model.PaymentProviderOpenPayment,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := topUp.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 创建充值订单失败 user_id=%d trade_no=%s error=%q", id, tradeNo, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	callBackAddress := service.GetCallbackAddress()
	notifyUrl := callBackAddress + "/api/user/open-payment/notify"
	returnUrl := setting.OpenPaymentReturnUrl
	if returnUrl == "" {
		returnUrl = paymentReturnPath("/console/topup")
	}

	orderReq := &paymentopen.CreateOrderRequest{
		ChannelCode:   method.ChannelCode,
		TradeType:     method.TradeType,
		Amount:        payMoney,
		Subject:       fmt.Sprintf("TUC%.2f", req.Amount),
		AppOrderNo:    tradeNo,
		NotifyUrl:     notifyUrl,
		ReturnUrl:     returnUrl,
		ExpireMinutes: setting.OpenPaymentOrderExpireMinutes,
	}

	orderResp, err := paymentopen.CreateOrder(orderReq)
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 下单失败 user_id=%d trade_no=%s payment_method=%s channel_code=%s trade_type=%s amount=%.2f notify_url=%s return_url=%s error=%q", id, tradeNo, req.PaymentMethod, method.ChannelCode, method.TradeType, payMoney, notifyUrl, returnUrl, err.Error()))
		c.JSON(http.StatusOK, gin.H{
			"message": "error",
			"data":    fmt.Sprintf("拉起支付失败：%s", err.Error()),
		})
		return
	}

	topUp.ProviderTradeNo = orderResp.OrderNo
	_ = topUp.Update()

	payMode := "redirect"
	if orderResp.ChannelResult.QrCodeUrl != "" {
		payMode = "qr"
	} else if len(orderResp.ChannelResult.PayParams) > 0 {
		payMode = "params"
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 订单创建成功 user_id=%d trade_no=%s order_no=%s pay_mode=%s", id, tradeNo, orderResp.OrderNo, payMode))

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"trade_no":    tradeNo,
			"order_no":    orderResp.OrderNo,
			"pay_mode":    payMode,
			"qr_code_url": orderResp.ChannelResult.QrCodeUrl,
			"pay_url":     orderResp.ChannelResult.PayUrl,
			"pay_params":  orderResp.ChannelResult.PayParams,
		},
	})
}

func OpenPaymentNotify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.LogWarn(c.Request.Context(), "支付中心 回调读取 body 失败")
		c.Status(http.StatusBadRequest)
		return
	}

	timestamp := c.GetHeader("X-Timestamp")
	nonce := c.GetHeader("X-Nonce")
	signature := c.GetHeader("X-Signature")

	if !paymentopen.VerifyNotification(setting.OpenPaymentAppSecret, body, timestamp, nonce, signature) {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("支付中心 回调验签失败 client_ip=%s", c.ClientIP()))
		c.Status(http.StatusForbidden)
		return
	}

	var payload paymentopen.NotifyPayload
	if err := common.Unmarshal(body, &payload); err != nil {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("支付中心 回调解析 JSON 失败 error=%q", err.Error()))
		c.Status(http.StatusOK)
		return
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 收到回调 event=%s app_order_no=%s order_no=%s", payload.Event, payload.AppOrderNo, payload.OrderNo))

	switch payload.Event {
	case "pay.success":
		handleOpenPaymentSuccess(c, &payload)
	case "pay.closed":
		handleOpenPaymentClosed(c, &payload)
	case "refund.success":
		handleOpenRefundSuccess(c, &payload)
	case "refund.failed":
		handleOpenRefundFailed(c, &payload)
	default:
		logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 忽略事件 event=%s", payload.Event))
	}

	c.Status(http.StatusOK)
}

func handleOpenPaymentSuccess(c *gin.Context, payload *paymentopen.NotifyPayload) {
	tradeNo := payload.AppOrderNo
	if tradeNo == "" {
		logger.LogWarn(c.Request.Context(), "支付中心 pay.success 缺少 app_order_no")
		return
	}

	LockOrder(tradeNo)
	defer UnlockOrder(tradeNo)

	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("支付中心 回调订单不存在 trade_no=%s", tradeNo))
		return
	}

	if topUp.PaymentProvider != model.PaymentProviderOpenPayment {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("支付中心 订单支付网关不匹配 trade_no=%s provider=%s", tradeNo, topUp.PaymentProvider))
		return
	}

	if topUp.Status == common.TopUpStatusSuccess {
		return
	}

	if topUp.Status != common.TopUpStatusPending {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("支付中心 订单状态异常 trade_no=%s status=%s", tradeNo, topUp.Status))
		return
	}

	callbackAmountYuan := decimal.NewFromInt(int64(payload.Amount)).Div(decimal.NewFromInt(100))
	localMoney := decimal.NewFromFloat(topUp.Money)
	if !callbackAmountYuan.Equal(localMoney) {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 回调金额不一致 trade_no=%s callback=%s local=%s", tradeNo, callbackAmountYuan.String(), localMoney.String()))
		return
	}

	topUp.Status = common.TopUpStatusSuccess
	topUp.CompleteTime = common.GetTimestamp()
	if err := topUp.Update(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 更新订单失败 trade_no=%s error=%q", tradeNo, err.Error()))
		return
	}

	dAmount := decimal.NewFromInt(int64(topUp.Amount))
	dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
	quotaToAdd := int(dAmount.Mul(dQuotaPerUnit).IntPart())
	if err := model.IncreaseUserQuota(topUp.UserId, quotaToAdd, true); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 更新用户额度失败 trade_no=%s user_id=%d quota=%d error=%q", tradeNo, topUp.UserId, quotaToAdd, err.Error()))
		return
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 充值成功 trade_no=%s user_id=%d quota=%d money=%.2f", tradeNo, topUp.UserId, quotaToAdd, topUp.Money))
	model.RecordTopupLog(topUp.UserId, fmt.Sprintf("使用在线充值成功，充值金额: %v，支付金额：%f", logger.LogQuota(quotaToAdd), topUp.Money), c.ClientIP(), topUp.PaymentMethod, "open_payment")
}

func handleOpenPaymentClosed(c *gin.Context, payload *paymentopen.NotifyPayload) {
	tradeNo := payload.AppOrderNo
	if tradeNo == "" {
		return
	}
	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil || topUp.Status != common.TopUpStatusPending {
		return
	}
	topUp.Status = common.TopUpStatusExpired
	_ = topUp.Update()
	logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 订单关闭 trade_no=%s", tradeNo))
}

func GetOpenPaymentStatus(c *gin.Context) {
	tradeNo := c.Query("trade_no")
	if tradeNo == "" {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "缺少 trade_no"})
		return
	}

	userId := c.GetInt("id")
	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil || topUp.UserId != userId {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "订单不存在"})
		return
	}

	if topUp.Status == common.TopUpStatusPending && topUp.ProviderTradeNo != "" {
		if resp, err := paymentopen.QueryOrder(topUp.ProviderTradeNo); err == nil && resp.Status == "paid" {
			handleOpenPaymentSuccess(c, &paymentopen.NotifyPayload{
				Event:      "pay.success",
				AppOrderNo: tradeNo,
				Amount:     resp.Amount,
			})
			topUp = model.GetTopUpByTradeNo(tradeNo)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"status":        topUp.Status,
			"trade_no":      topUp.TradeNo,
			"refund_status": topUp.RefundStatus,
		},
	})
}

type OpenPaymentRefundRequest struct {
	TradeNo      string  `json:"trade_no" binding:"required"`
	RefundAmount float64 `json:"refund_amount" binding:"required,gt=0"`
	Reason       string  `json:"reason"`
}

func RequestOpenPaymentRefund(c *gin.Context) {
	var req OpenPaymentRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	topUp := model.GetTopUpByTradeNo(req.TradeNo)
	if topUp == nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "订单不存在"})
		return
	}

	if topUp.PaymentProvider != model.PaymentProviderOpenPayment {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "仅支持支付中心订单退款"})
		return
	}

	if topUp.Status != common.TopUpStatusSuccess {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "订单未支付成功，无法退款"})
		return
	}

	maxRefundable := decimal.NewFromFloat(topUp.Money).Sub(decimal.NewFromFloat(topUp.RefundAmount))
	refundAmount := decimal.NewFromFloat(req.RefundAmount)
	if refundAmount.GreaterThan(maxRefundable) {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("退款金额不能超过可退金额 %s 元", maxRefundable.String())})
		return
	}

	dAmount := decimal.NewFromInt(int64(topUp.Amount))
	dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
	refundRatio := refundAmount.Div(decimal.NewFromFloat(topUp.Money))
	quotaToDeduct := int(dAmount.Mul(dQuotaPerUnit).Mul(refundRatio).IntPart())

	operatorId := c.GetInt("id")
	reason := req.Reason
	if reason == "" {
		reason = "管理员退款"
	}

	refundRecord := &model.TopUpRefund{
		TopUpId:         topUp.Id,
		UserId:          topUp.UserId,
		TradeNo:         topUp.TradeNo,
		ProviderOrderNo: topUp.ProviderTradeNo,
		RefundAmount:    req.RefundAmount,
		TotalAmount:     topUp.Money,
		Reason:          reason,
		Status:          model.RefundStatusPending,
		QuotaDeducted:   quotaToDeduct,
		OperatorId:      operatorId,
		CreateTime:      time.Now().Unix(),
	}

	if err := refundRecord.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 创建退款记录失败 trade_no=%s error=%q", req.TradeNo, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "创建退款记录失败"})
		return
	}

	refundReq := &paymentopen.CreateRefundRequest{
		OrderNo:      topUp.ProviderTradeNo,
		RefundAmount: req.RefundAmount,
		Reason:       reason,
	}

	refundResp, err := paymentopen.CreateRefund(refundReq)
	if err != nil {
		refundRecord.Status = model.RefundStatusFailed
		_ = refundRecord.Update()
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 退款请求失败 trade_no=%s error=%q", req.TradeNo, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "退款请求失败: " + err.Error()})
		return
	}

	refundRecord.RefundNo = refundResp.RefundNo
	refundRecord.Status = model.RefundStatusProcessing

	if refundResp.Status == "success" {
		refundRecord.Status = model.RefundStatusSuccess
		refundRecord.CompleteTime = time.Now().Unix()
		_ = refundRecord.Update()

		topUp.RefundAmount += req.RefundAmount
		if topUp.RefundAmount >= topUp.Money {
			topUp.RefundStatus = "refunded"
		} else {
			topUp.RefundStatus = "refunding"
		}
		_ = topUp.Update()

		if quotaToDeduct > 0 {
			_ = model.DecreaseUserQuota(topUp.UserId, quotaToDeduct, true)
		}
		model.RecordTopupLog(topUp.UserId, fmt.Sprintf("管理员退款成功，退款金额: %.2f 元，扣回额度: %v", req.RefundAmount, logger.LogQuota(quotaToDeduct)), c.ClientIP(), topUp.PaymentMethod, "open_payment_refund")
	} else {
		_ = refundRecord.Update()
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 退款发起成功 trade_no=%s refund_no=%s status=%s", req.TradeNo, refundResp.RefundNo, refundResp.Status))
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"refund_no": refundResp.RefundNo,
			"status":    refundRecord.Status,
		},
	})
}

func GetOpenPaymentRefunds(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	refunds, total := model.GetAllOpenPaymentRefunds(page, pageSize)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"list":  refunds,
			"total": total,
		},
	})
}

func handleOpenRefundSuccess(c *gin.Context, payload *paymentopen.NotifyPayload) {
	refundNo := payload.RefundNo
	if refundNo == "" {
		logger.LogWarn(c.Request.Context(), "支付中心 refund.success 缺少 refund_no")
		return
	}

	LockOrder(refundNo)
	defer UnlockOrder(refundNo)

	if err := model.CompleteOpenPaymentRefund(refundNo); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 处理退款成功回调失败 refund_no=%s error=%q", refundNo, err.Error()))
		return
	}

	refund := model.GetTopUpRefundByRefundNo(refundNo)
	if refund != nil {
		model.RecordTopupLog(refund.UserId, fmt.Sprintf("退款成功（回调），退款金额: %.2f 元，扣回额度: %v", refund.RefundAmount, logger.LogQuota(refund.QuotaDeducted)), c.ClientIP(), "", "open_payment_refund")
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("支付中心 退款成功 refund_no=%s", refundNo))
}

func handleOpenRefundFailed(c *gin.Context, payload *paymentopen.NotifyPayload) {
	refundNo := payload.RefundNo
	if refundNo == "" {
		return
	}

	LockOrder(refundNo)
	defer UnlockOrder(refundNo)

	if err := model.FailOpenPaymentRefund(refundNo); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("支付中心 处理退款失败回调失败 refund_no=%s error=%q", refundNo, err.Error()))
	}
	logger.LogWarn(c.Request.Context(), fmt.Sprintf("支付中心 退款失败 refund_no=%s reason=%s", refundNo, payload.Reason))
}
