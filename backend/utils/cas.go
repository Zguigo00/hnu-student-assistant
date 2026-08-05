package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// CASResult CAS 登录结果。
type CASResult struct {
	Client *resty.Client
	Ticket string
}

// CASLogin 执行通用 CAS 登录流程：
// 1. 获取登录页 → 提取隐藏字段 + AES 加密参数
// 2. 加密密码 → 提交表单（禁用重定向）
// 3. 从 302 Location 中提取 ticket
//
// tokenValidator 由调用方提供，负责用 ticket 换取 token。
func CASLogin(
	casLoginURL string,
	username, password string,
	tokenValidator func(client *resty.Client, ticket string) (string, error),
) (*resty.Client, error) {
	client := resty.New()
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// 获取 CAS 登录页
	loginResp, err := client.R().Get(casLoginURL)
	if err != nil {
		return nil, fmt.Errorf("访问 CAS 登录页失败: %w", err)
	}
	html := loginResp.String()

	// 提取隐藏字段 + AES 加密密码
	hiddenFields := ExtractHiddenFields(html)

	cp, err := ExtractCryptoParams(client, html)
	if err != nil {
		return nil, fmt.Errorf("提取加密参数失败: %w", err)
	}

	encryptedPwd, err := AesCBCEncrypt(password, cp.Key, cp.Iv)
	if err != nil {
		return nil, fmt.Errorf("AES 加密密码失败: %w", err)
	}

	// 构建表单
	formData := map[string]string{
		"username":  username,
		"password":  encryptedPwd,
		"lt":        hiddenFields["lt"],
		"execution": hiddenFields["execution"],
		"_eventId":  hiddenFields["_eventId"],
		"dllt":      hiddenFields["dllt"],
		"rmShown":   hiddenFields["rmShown"],
	}
	for k, v := range hiddenFields {
		if _, exists := formData[k]; !exists {
			formData[k] = v
		}
	}

	// 禁用重定向，提交登录表单
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	loginPostResp, err := client.R().
		SetFormData(formData).
		Post(casLoginURL)

	// 从 302 Location 中提取 ticket
	var ticket string
	if loginPostResp != nil {
		ticket = ExtractTicket(loginPostResp.Header().Get("Location"))
	}
	if ticket == "" {
		if loginPostResp != nil {
			body := loginPostResp.String()
			if strings.Contains(body, "密码有误") || strings.Contains(body, "用户名不存在") {
				return nil, fmt.Errorf("登录失败，请检查账号密码")
			}
		}
		if err != nil {
			return nil, fmt.Errorf("登录失败: %w", err)
		}
		return nil, fmt.Errorf("登录失败，未获取到 ticket")
	}

	// 用 ticket 换取 token
	token, err := tokenValidator(client, ticket)
	if err != nil {
		return nil, fmt.Errorf("获取 token 失败: %w", err)
	}

	// 恢复重定向，设置 token
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))
	client.SetHeader("X-Access-Token", token)

	return client, nil
}

// ExtractTicket 从重定向 URL 中提取 CAS ticket。
func ExtractTicket(location string) string {
	if location == "" {
		return ""
	}
	u, err := url.Parse(location)
	if err != nil {
		return ""
	}
	return u.Query().Get("ticket")
}

// ValidateCASTicket 通用 CAS ticket → JWT token 换取。
func ValidateCASTicket(validateURL, service string) func(client *resty.Client, ticket string) (string, error) {
	return func(client *resty.Client, ticket string) (string, error) {
		var result struct {
			Success bool `json:"success"`
			Result  struct {
				Token string `json:"token"`
			} `json:"result"`
		}

		resp, err := client.R().
			SetQueryParam("_t", fmt.Sprintf("%d", time.Now().UnixMilli())).
			SetQueryParam("ticket", ticket).
			SetQueryParam("service", service).
			Get(validateURL)
		if err != nil {
			return "", fmt.Errorf("验证 ticket 失败: %w", err)
		}

		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return "", fmt.Errorf("解析 token 响应失败: %w", err)
		}
		if !result.Success || result.Result.Token == "" {
			return "", fmt.Errorf("获取 token 失败: %s", string(resp.Body()))
		}

		return result.Result.Token, nil
	}
}
