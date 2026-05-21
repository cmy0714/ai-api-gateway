package setting

import (
	"strings"
)

var OpenPaymentBaseURL = ""
var OpenPaymentAppId = ""
var OpenPaymentAppSecret = ""
var OpenPaymentReturnUrl = ""
var OpenPaymentMinTopUp = 1.0
var OpenPaymentOrderExpireMinutes = 30
var OpenPaymentMethods = "[]"

type OpenPaymentMethod struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ChannelCode string `json:"channel_code"`
	TradeType   string `json:"trade_type"`
	Color       string `json:"color"`
}

func OpenPaymentMethodFromPaymentType(paymentMethod string) (OpenPaymentMethod, bool) {
	paymentType := NormalizeOpenPaymentMethodType(paymentMethod)
	switch paymentType {
	case "open_alipay":
		return OpenPaymentMethod{
			Name:        "支付宝",
			Type:        "alipay",
			ChannelCode: "alipay",
			TradeType:   "page",
			Color:       "#1677FF",
		}, true
	case "open_wxpay", "open_wechat", "open_wechat_native":
		return OpenPaymentMethod{
			Name:        "微信支付",
			Type:        "wxpay",
			ChannelCode: "wechat",
			TradeType:   "native",
			Color:       "#07C160",
		}, true
	default:
		return OpenPaymentMethod{}, false
	}
}

func OpenPaymentMethodFromPayMethod(payMethod map[string]string) (OpenPaymentMethod, bool) {
	openMethod, ok := OpenPaymentMethodFromPaymentType(payMethod["type"])
	if !ok {
		name := strings.ToLower(strings.TrimSpace(payMethod["name"]))
		switch {
		case strings.Contains(name, "支付宝") || strings.Contains(name, "alipay"):
			openMethod, ok = OpenPaymentMethodFromPaymentType("alipay")
		case strings.Contains(name, "微信") || strings.Contains(name, "wechat") || strings.Contains(name, "weixin"):
			openMethod, ok = OpenPaymentMethodFromPaymentType("wxpay")
		}
	}
	if !ok {
		return OpenPaymentMethod{}, false
	}
	if payMethod["name"] != "" {
		openMethod.Name = payMethod["name"]
	}
	if payMethod["color"] != "" {
		openMethod.Color = payMethod["color"]
	}
	return openMethod, true
}

func GetOpenPaymentMethodsFromPayMethods(payMethods []map[string]string) []OpenPaymentMethod {
	methods := make([]OpenPaymentMethod, 0, len(payMethods))
	for _, method := range payMethods {
		openMethod, ok := OpenPaymentMethodFromPayMethod(method)
		if !ok {
			continue
		}
		methods = append(methods, openMethod)
	}
	return methods
}

func GetOpenPaymentMethodsForAPI(payMethods []map[string]string) []map[string]string {
	methods := GetOpenPaymentMethodsFromPayMethods(payMethods)
	result := make([]map[string]string, 0, len(methods))
	for _, m := range methods {
		result = append(result, map[string]string{
			"name":     m.Name,
			"type":     NormalizeOpenPaymentMethodType(m.Type),
			"color":    m.Color,
			"provider": "open_payment",
		})
	}
	return result
}

func FindOpenPaymentMethod(paymentMethod string, payMethods []map[string]string) *OpenPaymentMethod {
	methods := GetOpenPaymentMethodsFromPayMethods(payMethods)
	for _, m := range methods {
		if m.Type == paymentMethod || NormalizeOpenPaymentMethodType(m.Type) == paymentMethod {
			return &m
		}
	}
	return nil
}

func NormalizeOpenPaymentMethodType(paymentMethod string) string {
	paymentMethod = strings.ToLower(strings.TrimSpace(paymentMethod))
	if strings.HasPrefix(paymentMethod, "open_") {
		return paymentMethod
	}
	return "open_" + paymentMethod
}
