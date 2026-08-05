package portal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// NewsService 校园新闻查询服务。
type NewsService struct {
	httpClient *resty.Client
}

func NewNewsService(client *resty.Client) *NewsService {
	return &NewsService{httpClient: client}
}

// GetNews 获取校园新闻。startIndex/endIndex 控制分页范围。
func (s *NewsService) GetNews(startIndex, endIndex int) ([]NewsItem, int, error) {
	resp, err := s.httpClient.R().
		SetQueryParam("_t", fmt.Sprintf("%d", time.Now().UnixMilli())).
		SetQueryParam("startIndex", fmt.Sprintf("%d", startIndex)).
		SetQueryParam("endIndex", fmt.Sprintf("%d", endIndex)).
		Get(NewsAPI)
	if err != nil {
		return nil, 0, fmt.Errorf("请求校园新闻失败: %w", err)
	}

	// 先解析外层结构，result 可能是对象或数组
	var raw struct {
		Success bool            `json:"success"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		return nil, 0, fmt.Errorf("解析新闻响应失败: %w", err)
	}
	if !raw.Success {
		return nil, 0, fmt.Errorf("新闻接口返回失败 (HTTP %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	// 尝试 result.records 格式
	var withRecords struct {
		Records []NewsItem `json:"records"`
	}
	if err := json.Unmarshal(raw.Result, &withRecords); err == nil && len(withRecords.Records) > 0 {
		return withRecords.Records, len(withRecords.Records), nil
	}

	// 尝试 result 直接是数组的格式
	var items []NewsItem
	if err := json.Unmarshal(raw.Result, &items); err == nil {
		return items, len(items), nil
	}

	return nil, 0, fmt.Errorf("无法解析新闻数据: %s", string(raw.Result))
}
