package portal

import "encoding/json"

// NewsItem 新闻记录
type NewsItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime string `json:"createTime"`
	CreateBy   string `json:"createBy"`
	Look       int    `json:"look"`
	Content    string `json:"content"`
	IsLinks    bool   `json:"isLinks"`
	Type       int    `json:"type"`
}

// UnmarshalJSON 兼容 isLinks 为 bool 或 number 的两种 API 格式
func (n *NewsItem) UnmarshalJSON(data []byte) error {
	type Alias NewsItem
	aux := &struct {
		IsLinks json.RawMessage `json:"isLinks"`
		*Alias
	}{Alias: (*Alias)(n)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	n.IsLinks = false
	if len(aux.IsLinks) > 0 {
		s := string(aux.IsLinks)
		if s == "true" || s == "1" {
			n.IsLinks = true
		}
	}

	return nil
}
