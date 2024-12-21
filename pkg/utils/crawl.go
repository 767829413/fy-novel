package utils

import (
	"net/http"
	"net/url"
	"strings"
)

// NormalizeURL 将相对路径转换为绝对路径，并确保 URL 格式正确
func NormalizeURL(s, host string) string {

	// 检查 s 是否已经是一个完整的 URL
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return cleanURL(s)
	}

	// 如果 s 不是完整 URL，则将其与 host 组合
	baseURL, err := url.Parse(host)
	if err != nil {
		return ""
	}

	relativeURL, err := url.Parse(s)
	if err != nil {
		return ""
	}

	return cleanURL(baseURL.ResolveReference(relativeURL).String())
}

// cleanURL 清理 URL，移除多余的斜杠和处理编码
func cleanURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	// 确保路径不以 '//' 开始
	u.Path = strings.TrimPrefix(u.Path, "//")

	// 移除路径中的连续斜杠
	for strings.Contains(u.Path, "//") {
		u.Path = strings.ReplaceAll(u.Path, "//", "/")
	}

	// 确保查询参数正确编码
	q := u.Query()
	u.RawQuery = q.Encode()

	return u.String()
}

// BuildParams builds a map of parameters from a JSON string and a keyword
func BuildParams(body map[string]string, keyword, searchField string) map[string]string {
	if len(body) == 0 {
		return nil
	}
	if keyword == "" {
		return nil
	}

	params := make(map[string]string)

	for key, value := range body {
		if key == searchField {
			params[value] = keyword
		} else {
			params[key] = value
		}
	}

	return params
}

// BuildMethod returns the corresponding HTTP method
func BuildMethod(method string) string {
	if method == "" {
		return http.MethodPost // 默认返回 POST 方法
	}

	switch strings.ToLower(method) {
	case "get":
		return http.MethodGet
	case "post":
		return http.MethodPost
	case "put":
		return http.MethodPut
	case "delete":
		return http.MethodDelete
	case "patch":
		return http.MethodPatch
	case "head":
		return http.MethodHead
	case "options":
		return http.MethodOptions
	case "trace":
		return http.MethodTrace
	default:
		return http.MethodPost
	}
}
