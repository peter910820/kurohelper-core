package kurohelpercore

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
