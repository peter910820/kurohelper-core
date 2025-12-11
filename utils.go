package kurohelpercore

import "net/url"

// check if the string is English
func IsEnglish(r rune) bool {
	if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
		return false
	}
	return true
}

// 從any類型中提取字串值(安全行為)
func GetStringValue(value any) string {
	if str, ok := value.(string); ok {
		return str
	}
	return "wrong type"
}

// 檢查URL是否合法
func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// scheme 必須是 http 或 https
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// host 不能為空
	if u.Host == "" {
		return false
	}

	return true
}
