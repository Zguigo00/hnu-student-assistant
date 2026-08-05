package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseXnxq 将学年学期字符串（如 "2024-2025-1"）转换为教务系统所需的 xnm、xqm。
// xqm 映射："1" → "3"（第一学期），"2" → "12"（第二学期）。
func ParseXnxq(xnxq string) (xnm, xqm string) {
	if xnxq == "" {
		return "", ""
	}
	parts := strings.Split(xnxq, "-")
	if len(parts) >= 2 {
		xnm = parts[0]
	}
	if len(parts) >= 3 {
		switch parts[2] {
		case "1":
			xqm = "3"
		case "2":
			xqm = "12"
		}
	}
	return
}

// ParseFloat 安全地将字符串转为 float64，解析失败返回 0。
func ParseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// Stringify 将 interface{} 安全地转为字符串。
func Stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

var (
	hiddenInputRe = regexp.MustCompile(`(?i)<input[^>]*type="hidden"[^>]*>`)
	hiddenNameRe  = regexp.MustCompile(`(?i)name="([^"]*)"`)
	hiddenValueRe = regexp.MustCompile(`(?i)value\s*=\s*(?:"([^"]*)"|([^>\s]+))`)
)

// ExtractHiddenFields 从 HTML 中提取所有隐藏表单字段。
func ExtractHiddenFields(html string) map[string]string {
	fields := make(map[string]string)
	for _, input := range hiddenInputRe.FindAllString(html, -1) {
		nameMatch := hiddenNameRe.FindStringSubmatch(input)
		if len(nameMatch) <= 1 {
			continue
		}
		value := ""
		if valueMatch := hiddenValueRe.FindStringSubmatch(input); len(valueMatch) > 1 {
			if valueMatch[1] != "" {
				value = valueMatch[1]
			} else if len(valueMatch) > 2 {
				value = valueMatch[2]
			}
		}
		fields[nameMatch[1]] = value
	}
	return fields
}
