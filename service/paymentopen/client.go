package paymentopen

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func normalizeBaseURL(rawURL string) string {
	u := strings.TrimRight(rawURL, "/")
	for _, suffix := range []string{"/open/payment", "/open/payment/orders"} {
		if strings.HasSuffix(u, suffix) {
			u = u[:len(u)-len(suffix)]
			break
		}
	}
	return u
}

func apiBase() string {
	baseURL := normalizeBaseURL(setting.OpenPaymentBaseURL)
	for _, suffix := range []string{"/auth-app-api/v1", "/admin-api/v1", "/api/v1"} {
		if strings.HasSuffix(baseURL, suffix) {
			return baseURL + "/open/payment"
		}
	}
	return baseURL + "/api/v1/open/payment"
}

func doRequest(method, url string, body []byte) ([]byte, error) {
	headers := buildSignedHeaders(setting.OpenPaymentAppId, setting.OpenPaymentAppSecret, method, url, body)

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d %s: %s", resp.StatusCode, url, string(respBody))
	}

	return respBody, nil
}

func parseResponse(respBody []byte) (*APIResponse, error) {
	var apiResp APIResponse
	if err := common.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	if apiResp.Code != 1 {
		return nil, fmt.Errorf("payment center error [%d]: %s", apiResp.Code, apiResp.Msg)
	}
	return &apiResp, nil
}

// CreateOrder creates a payment order at the payment center.
func CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
	url := apiBase() + "/orders"
	body, err := common.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBody, err := doRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	apiResp, err := parseResponse(respBody)
	if err != nil {
		return nil, err
	}

	dataBytes, err := common.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var result CreateOrderResponse
	if err := common.Unmarshal(dataBytes, &result); err != nil {
		return nil, err
	}

	if result.Order != nil && result.OrderNo == "" {
		result.OrderNo = result.Order.OrderNo
	}

	return &result, nil
}

// QueryOrder queries payment order status from the payment center.
func QueryOrder(orderNo string) (*QueryOrderResponse, error) {
	url := apiBase() + "/orders/" + orderNo

	respBody, err := doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	apiResp, err := parseResponse(respBody)
	if err != nil {
		return nil, err
	}

	dataBytes, err := common.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var result QueryOrderResponse
	if err := common.Unmarshal(dataBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CloseOrder closes an unpaid order.
func CloseOrder(orderNo string) error {
	url := apiBase() + "/orders/" + orderNo + "/close"

	respBody, err := doRequest("POST", url, nil)
	if err != nil {
		return err
	}

	_, err = parseResponse(respBody)
	return err
}

// CreateRefund initiates a refund through the payment center.
func CreateRefund(req *CreateRefundRequest) (*CreateRefundResponse, error) {
	url := apiBase() + "/refunds"
	body, err := common.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBody, err := doRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	apiResp, err := parseResponse(respBody)
	if err != nil {
		return nil, err
	}

	dataBytes, err := common.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var result CreateRefundResponse
	if err := common.Unmarshal(dataBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// IsConfigured returns true if all required open payment settings are set.
func IsConfigured() bool {
	return strings.TrimSpace(setting.OpenPaymentBaseURL) != "" &&
		strings.TrimSpace(setting.OpenPaymentAppId) != "" &&
		strings.TrimSpace(setting.OpenPaymentAppSecret) != ""
}
