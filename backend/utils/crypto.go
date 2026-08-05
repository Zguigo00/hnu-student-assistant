package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
)

// CryptoParams AES 加密参数。
type CryptoParams struct {
	Key []byte
	Iv  []byte
}

// ExtractCryptoParams 从 CAS 登录页中提取 AES 加密参数。
func ExtractCryptoParams(client *resty.Client, html string) (CryptoParams, error) {
	allContent := html

	jsURLs := FindCryptoJSURLs(html)
	for _, jsURL := range jsURLs {
		jsContent, err := FetchJS(client, jsURL)
		if err != nil {
			continue
		}
		allContent += "\n" + jsContent
	}

	return ExtractFromHTML(allContent)
}

// FindCryptoJSURLs 从 HTML 中查找加密相关的 JS 文件 URL。
func FindCryptoJSURLs(html string) []string {
	var urls []string
	seen := map[string]bool{}

	addURL := func(url string) {
		if strings.HasPrefix(url, "/") {
			url = "https://auth.chnu.edu.cn" + url
		} else if !strings.HasPrefix(url, "http") {
			url = "https://auth.chnu.edu.cn/authserver/" + url
		}
		if !seen[url] {
			seen[url] = true
			urls = append(urls, url)
		}
	}

	re := regexp.MustCompile(`(?i)<script[^>]+src="([^"]+\.js[^"]*)"`)
	for _, m := range re.FindAllStringSubmatch(html, -1) {
		if len(m) > 1 {
			addURL(m[1])
		}
	}

	re2 := regexp.MustCompile(`loadJavascript\s*\(\s*(?:baseName\s*\+\s*)?['"]([^'"]+)['"]`)
	for _, m := range re2.FindAllStringSubmatch(html, -1) {
		if len(m) > 1 {
			addURL(m[1])
		}
	}

	return urls
}

// FetchJS 获取外部 JS 文件内容。
func FetchJS(client *resty.Client, url string) (string, error) {
	resp, err := client.R().Get(url)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

// ExtractFromHTML 从 HTML/JS 内容中提取加密参数。
func ExtractFromHTML(content string) (CryptoParams, error) {
	encLitRe := regexp.MustCompile(`encryption\s*\([^,]+,\s*['"]([^'"]+)['"][,\s]*['"]([^'"]+)['"]`)
	if matches := encLitRe.FindStringSubmatch(content); len(matches) > 2 {
		return CryptoParams{Key: []byte(matches[1]), Iv: []byte(matches[2])}, nil
	}

	encVarRe := regexp.MustCompile(`encryption\s*\([^,]+,\s*(\w+)\s*,\s*(\w+)\s*\)`)
	if matches := encVarRe.FindStringSubmatch(content); len(matches) > 2 {
		keyVal := FindScriptVar(content, matches[1])
		ivVal := FindScriptVar(content, matches[2])
		if keyVal != "" && ivVal != "" {
			return CryptoParams{Key: []byte(keyVal), Iv: []byte(ivVal)}, nil
		}
	}

	re := regexp.MustCompile(`pwdDefaultEncryptSalt\s*=\s*"([^"]+)"`)
	if matches := re.FindStringSubmatch(content); len(matches) > 1 {
		salt := matches[1]
		key := []byte(salt)
		if len(key) > 16 {
			key = key[:16]
		}
		iv := make([]byte, 16)
		copy(iv, key)
		return CryptoParams{Key: key, Iv: iv}, nil
	}

	return CryptoParams{}, fmt.Errorf("未找到 AES 加密参数")
}

// FindScriptVar 从 HTML 的 <script> 中查找 JavaScript 变量的字符串值。
func FindScriptVar(html, varName string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?:var\s+)?%s\s*=\s*['"]([^'"]+)['"]`, regexp.QuoteMeta(varName)))
	if matches := re.FindStringSubmatch(html); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// AesCBCEncrypt 使用 AES-CBC 模式加密，Zero Padding，返回 Base64 字符串。
func AesCBCEncrypt(plaintext string, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	blockSize := block.BlockSize()
	padding := blockSize - len(plaintext)%blockSize
	if padding == blockSize {
		padding = 0
	}
	padded := make([]byte, len(plaintext)+padding)
	copy(padded, plaintext)

	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
