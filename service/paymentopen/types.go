package paymentopen

type CreateOrderRequest struct {
	ChannelCode   string  `json:"channel_code"`
	TradeType     string  `json:"trade_type"`
	Amount        float64 `json:"amount"`
	Subject       string  `json:"subject"`
	AppOrderNo    string  `json:"app_order_no,omitempty"`
	NotifyUrl     string  `json:"notify_url,omitempty"`
	ReturnUrl     string  `json:"return_url,omitempty"`
	ExpireMinutes int     `json:"expire_minutes,omitempty"`
	Openid        string  `json:"openid,omitempty"`
}

type CreateOrderResponse struct {
	OrderNo       string            `json:"order_no"`
	AppOrderNo    string            `json:"app_order_no"`
	Amount        int               `json:"amount"`
	Status        string            `json:"status"`
	Idempotent    bool              `json:"idempotent"`
	ChannelResult ChannelResult     `json:"channel_result"`
	Order         *OrderInfo        `json:"order,omitempty"`
}

type ChannelResult struct {
	Success   bool              `json:"success"`
	QrCodeUrl string            `json:"qr_code_url"`
	PayUrl    string            `json:"pay_url"`
	PayParams map[string]interface{} `json:"pay_params"`
	ErrorMsg  string            `json:"error_msg"`
}

type OrderInfo struct {
	OrderNo   string `json:"order_no"`
	ExpiredAt string `json:"expired_at"`
}

type QueryOrderResponse struct {
	OrderNo        string `json:"order_no"`
	Status         string `json:"status"`
	ChannelOrderNo string `json:"channel_order_no"`
	PaidAt         string `json:"paid_at"`
	Amount         int    `json:"amount"`
	RefundAmount   int    `json:"refund_amount"`
}

type CreateRefundRequest struct {
	OrderNo      string  `json:"order_no"`
	RefundAmount float64 `json:"refund_amount"`
	Reason       string  `json:"reason"`
}

type CreateRefundResponse struct {
	RefundNo string `json:"refund_no"`
	Status   string `json:"status"`
}

type APIResponse struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data interface{}      `json:"data"`
}

type NotifyPayload struct {
	Event          string `json:"event"`
	OrderNo        string `json:"order_no"`
	AppOrderNo     string `json:"app_order_no"`
	ChannelCode    string `json:"channel_code"`
	ChannelOrderNo string `json:"channel_order_no"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	Subject        string `json:"subject"`
	Status         string `json:"status"`
	PaidAt         string `json:"paid_at"`
	ClosedAt       string `json:"closed_at"`
	RefundNo       string `json:"refund_no"`
	RefundAmount   int    `json:"refund_amount"`
	TotalAmount    int    `json:"total_amount"`
	Reason         string `json:"reason"`
	RefundedAt     string `json:"refunded_at"`
}
