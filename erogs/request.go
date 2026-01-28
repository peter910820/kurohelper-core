package erogs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	kurohelpercore "kurohelper-core"

	kurohelperproxy "github.com/kuro-helper/kurohelper-proxy"
)

type rateLimitStruct struct {
	Quota     int
	ResetTime time.Time
	RWMu      sync.RWMutex
}

var (
	resetTime       time.Duration
	rateLimitRecord rateLimitStruct
)

// 確保設定檔初始化後才初始化速率鎖的變數
func InitRateLimit(resetTime time.Duration) {
	resetTime = time.Duration(resetTime) * time.Second
	rateLimitRecord = rateLimitStruct{
		Quota:     5,
		ResetTime: time.Now().Add(resetTime),
	}
}

func sendPostRequest(sql string) (string, error) {
	if !rateLimit(1) {
		return "", kurohelpercore.ErrRateLimit
	}

	formData := url.Values{}
	formData.Set("sql", sql)

	req, err := http.NewRequest(http.MethodPost, os.Getenv("EROGS_ENDPOINT"), strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")

	var client *http.Client
	if os.Getenv("PROXY_USE") == "private" {
		dialer, err := kurohelperproxy.GetProxyDialer(os.Getenv("PROXY_PRIVATE_IP"), nil, os.Getenv("PROXY_PRIVATE_PORT"))
		if err != nil {
			fmt.Print(err)
			return "", err
		}
		client = &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
			},
			Timeout: 10 * time.Second,
		}
	} else if os.Getenv("PROXY_USE") == "vpn" {
		proxyAuth := kurohelperproxy.GenerateProxyAuth(os.Getenv("PROXY_AUTH_USER"), os.Getenv("PROXY_AUTH_PASSWORD"))
		dialer, err := kurohelperproxy.GetProxyDialer(os.Getenv("PROXY_ADDRESS"), proxyAuth, os.Getenv("PROXY_PORT"))
		if err != nil {
			fmt.Print(err)
			return "", err
		}
		client = &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
			},
			Timeout: 15 * time.Second,
		}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%w %d", kurohelpercore.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return "", err
	}

	selection := doc.Find("td").First() // 只取第一個符合的
	jsonText := selection.Text()

	if strings.TrimSpace(jsonText) == "" {
		return "", kurohelpercore.ErrSearchNoContent
	}

	return jsonText, nil
}

func rateLimit(quota int) bool {
	rateLimitRecord.RWMu.Lock()
	defer rateLimitRecord.RWMu.Unlock()

	now := time.Now()
	if now.After(rateLimitRecord.ResetTime) {
		rateLimitRecord.Quota = 5
		rateLimitRecord.ResetTime = now.Add(resetTime)
	}

	if rateLimitRecord.Quota > 0 {
		rateLimitRecord.Quota -= quota
		return true
	}
	return false
}
