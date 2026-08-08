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

// GetNewsByCategory 按分类获取校园新闻。
// category: "consulting"(新闻咨询) / "school"(学校新闻) / "notice"(通知公告) / "weekly"(周工作安排)
// page 从 1 开始，size 为每页条数。
func (s *NewsService) GetNewsByCategory(category string, page, size int) ([]NewsItem, int, error) {
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())
	req := s.httpClient.R().SetQueryParam("_t", ts)

	var url string
	switch category {
	case "school":
		url = ListNewsAPI
		req.SetQueryParam("newsType", "1").
			SetQueryParam("pageNo", fmt.Sprintf("%d", page)).
			SetQueryParam("pageSize", fmt.Sprintf("%d", size))
	case "weekly":
		url = ListNewsAPI
		req.SetQueryParam("newsType", "2").
			SetQueryParam("pageNo", fmt.Sprintf("%d", page)).
			SetQueryParam("pageSize", fmt.Sprintf("%d", size))
	case "notice":
		url = ListNoticeAPI
		req.SetQueryParam("pageNo", fmt.Sprintf("%d", page)).
			SetQueryParam("pageSize", fmt.Sprintf("%d", size))
	default: // consulting
		url = NewsAPI
		start := (page - 1) * size
		req.SetQueryParam("startIndex", fmt.Sprintf("%d", start)).
			SetQueryParam("endIndex", fmt.Sprintf("%d", start+size))
	}

	resp, err := req.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("请求新闻失败: %w", err)
	}

	return parseNewsResponse(resp)
}

func parseNewsResponse(resp *resty.Response) ([]NewsItem, int, error) {
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

	// 策略1: result.records 格式
	var withRecords struct {
		Records []NewsItem `json:"records"`
		Total   int        `json:"total"`
	}
	if err := json.Unmarshal(raw.Result, &withRecords); err == nil && len(withRecords.Records) > 0 {
		total := withRecords.Total
		if total == 0 {
			total = len(withRecords.Records)
		}
		return withRecords.Records, total, nil
	}

	// 策略2: result 直接是数组
	var items []NewsItem
	if err := json.Unmarshal(raw.Result, &items); err == nil && len(items) > 0 {
		return items, len(items), nil
	}

	return nil, 0, fmt.Errorf("无法解析新闻数据: %s", truncate(string(raw.Result), 500))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
