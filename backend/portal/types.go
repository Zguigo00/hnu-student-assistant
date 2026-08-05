package portal

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
