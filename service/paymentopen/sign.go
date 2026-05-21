package paymentopen

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func computeSignature(appSecret, timestamp, nonce, method, path, sortedQuery, bodyHash string) string {
	signString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", timestamp, nonce, method, path, sortedQuery, bodyHash)
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(signString))
	return hex.EncodeToString(mac.Sum(nil))
}

func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func buildSignedHeaders(appId, appSecret, method, rawURL string, body []byte) map[string]string {
	parsed, _ := url.Parse(rawURL)
	path := parsed.Path

	pairs := parsed.Query()
	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedParts := make([]string, 0, len(keys))
	for _, k := range keys {
		for _, v := range pairs[k] {
			sortedParts = append(sortedParts, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}
	sortedQuery := strings.Join(sortedParts, "&")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strings.ReplaceAll(uuid.New().String(), "-", "")
	bodyHash := sha256Hex(body)

	signature := computeSignature(appSecret, timestamp, nonce, strings.ToUpper(method), path, sortedQuery, bodyHash)

	return map[string]string{
		"X-App-Id":     appId,
		"X-Timestamp":  timestamp,
		"X-Nonce":      nonce,
		"X-Signature":  signature,
		"Content-Type": "application/json",
	}
}

// VerifyNotification verifies HMAC signature from payment center callback headers.
// Returns true if signature is valid.
func VerifyNotification(appSecret string, body []byte, timestamp, nonce, signature string) bool {
	if signature == "" || timestamp == "" || nonce == "" {
		return false
	}

	ts, err := fmt.Sscanf(timestamp, "%d", new(int64))
	if err != nil || ts != 1 {
		return false
	}

	var tsVal int64
	fmt.Sscanf(timestamp, "%d", &tsVal)
	now := time.Now().Unix()
	if now-tsVal > 300 || tsVal-now > 300 {
		return false
	}

	bodyHash := sha256Hex(body)
	signString := fmt.Sprintf("%s\n%s\n%s", timestamp, nonce, bodyHash)
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(signString))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}
