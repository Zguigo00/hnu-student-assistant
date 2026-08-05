package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"

	"hnu-student-assistant/backend/utils"

	"github.com/go-resty/resty/v2"
)

const (
	jwxtLoginPage = "https://jwglxt.chnu.edu.cn/xtgl/login_slogin.html"
	jwxtPublicKey = "https://jwglxt.chnu.edu.cn/xtgl/login_getPublicKey.html"
	jwxtLoginAPI  = "https://jwglxt.chnu.edu.cn/xtgl/login_slogin.html"
)

// JwxtLogin 登录教务系统（RSA 加密密码）。
// 返回已建立 session 的 resty.Client，供成绩/课表等服务使用。
func JwxtLogin(username, password string) (*resty.Client, error) {
	client := resty.New()
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	loginPageResp, err := client.R().Get(jwxtLoginPage)
	if err != nil {
		return nil, fmt.Errorf("访问教务系统登录页失败: %w", err)
	}

	csrftoken, err := extractCsrftoken(loginPageResp.String())
	if err != nil {
		return nil, fmt.Errorf("提取 csrftoken 失败: %w", err)
	}

	pubKeyResp, err := client.R().Get(jwxtPublicKey)
	if err != nil {
		return nil, fmt.Errorf("获取公钥失败: %w", err)
	}

	publicKey, err := parsePublicKey(pubKeyResp.String())
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("RSA 加密密码失败: %w", err)
	}
	encodedPwd := base64.StdEncoding.EncodeToString(encrypted)

	hiddenFields := utils.ExtractHiddenFields(loginPageResp.String())
	hiddenFields["yhm"] = username
	hiddenFields["mm"] = encodedPwd
	hiddenFields["csrftoken"] = csrftoken

	loginResp, err := client.R().
		SetFormData(hiddenFields).
		Post(jwxtLoginAPI)
	if err != nil {
		return nil, fmt.Errorf("登录请求失败: %w", err)
	}

	body := loginResp.String()
	if loginResp.StatusCode() != 200 {
		return nil, fmt.Errorf("登录失败，状态码: %d", loginResp.StatusCode())
	}
	if regexp.MustCompile(`name="csrftoken"`).MatchString(body) &&
		regexp.MustCompile(`name="mm"`).MatchString(body) {
		return nil, fmt.Errorf("登录失败，请检查学号和密码")
	}

	return client, nil
}

var csrftokenRe = regexp.MustCompile(`(?i)<input[^>]*name="csrftoken"[^>]*value\s*=\s*"([^"]*)"`)
var csrftokenRe2 = regexp.MustCompile(`(?i)<input[^>]*value\s*=\s*"([^"]*)"[^>]*name="csrftoken"`)

func extractCsrftoken(html string) (string, error) {
	matches := csrftokenRe.FindStringSubmatch(html)
	if len(matches) < 2 {
		matches = csrftokenRe2.FindStringSubmatch(html)
		if len(matches) < 2 {
			return "", fmt.Errorf("未在登录页中找到 csrftoken")
		}
	}
	return matches[1], nil
}

func parsePublicKey(body string) (*rsa.PublicKey, error) {
	var resp struct {
		Modulus  string `json:"modulus"`
		Exponent string `json:"exponent"`
	}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, fmt.Errorf("解析公钥 JSON 失败: %w", err)
	}

	modBytes, err := base64.StdEncoding.DecodeString(resp.Modulus)
	if err != nil {
		return nil, fmt.Errorf("解码 modulus 失败: %w", err)
	}

	expBytes, err := base64.StdEncoding.DecodeString(resp.Exponent)
	if err != nil {
		return nil, fmt.Errorf("解码 exponent 失败: %w", err)
	}
	expInt := 0
	for _, b := range expBytes {
		expInt = expInt<<8 | int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(modBytes),
		E: expInt,
	}, nil
}
