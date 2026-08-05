package sc

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const scoreAPI = "https://sc.chnu.edu.cn/hbsdDkApi/score/scScoreCollect/hour/statistics/list"

// ScoreService 第课堂成绩查询服务。
type ScoreService struct {
	httpClient *resty.Client
}

func NewScoreService(client *resty.Client) *ScoreService {
	return &ScoreService{httpClient: client}
}

// GetScores 查询第二课堂成绩。
// startTime/endTime 格式如 "2025-09-16 00:01:00"。
func (s *ScoreService) GetScores(startTime, endTime string) ([]SecondClassScore, error) {
	tableName := "userCode,userName,grade,college,major,classes,sxylyagrx,xskjycxcy,tydlyydjn,rwskyyssy,shsjyzyfw"
	field := "id,,userCode,userName,grade,college,major,classes,sxylyagrx,xskjycxcy,tydlyydjn,rwskyyssy,shsjyzyfw"

	resp, err := s.httpClient.R().
		SetQueryParam("tableName", tableName).
		SetQueryParam("startTime", startTime).
		SetQueryParam("endTime", endTime).
		SetQueryParam("column", "createTime").
		SetQueryParam("order", "desc").
		SetQueryParam("field", field).
		SetQueryParam("pageNo", "1").
		SetQueryParam("pageSize", "50").
		Get(scoreAPI)
	if err != nil {
		return nil, fmt.Errorf("查询第二课堂成绩失败: %w", err)
	}

	body := resp.String()

	// 先解析外层结构
	var raw struct {
		Success bool            `json:"success"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		return nil, fmt.Errorf("解析第二课堂响应失败: %w", err)
	}
	if !raw.Success {
		return nil, fmt.Errorf("第二课堂接口返回失败: %s", body)
	}

	// 尝试 result 是数组
	var items []SecondClassScore
	if err := json.Unmarshal(raw.Result, &items); err == nil && len(items) > 0 {
		return items, nil
	}

	// 尝试 result.records 格式（分页响应）
	var withRecords struct {
		Records []SecondClassScore `json:"records"`
	}
	if err := json.Unmarshal(raw.Result, &withRecords); err == nil && len(withRecords.Records) > 0 {
		return withRecords.Records, nil
	}

	// 尝试 result.data 格式
	var withData struct {
		Data []SecondClassScore `json:"data"`
	}
	if err := json.Unmarshal(raw.Result, &withData); err == nil && len(withData.Data) > 0 {
		return withData.Data, nil
	}

	return nil, fmt.Errorf("无法解析第二课堂数据: %s", string(raw.Result))
}
