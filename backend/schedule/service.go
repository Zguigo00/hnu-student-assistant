package schedule

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"hnu-student-assistant/backend/utils"

	"github.com/go-resty/resty/v2"
)

const scheduleURL = "https://jwglxt.chnu.edu.cn/kbcx/xskbcx_cxXsgrkb.html"

// ScheduleService 课表查询服务。
type ScheduleService struct {
	httpClient *resty.Client
}

func NewScheduleService(client *resty.Client) *ScheduleService {
	return &ScheduleService{httpClient: client}
}

// GetSchedule 查询课表。xnxq 格式如 "2024-2025-1"。
func (s *ScheduleService) GetSchedule(xnxq string) ([]ScheduleCourse, error) {
	xnm, xqm := utils.ParseXnxq(xnxq)

	resp, err := s.httpClient.R().
		SetQueryParam("gnmkdm", "N2151").
		SetFormData(map[string]string{
			"xnm": xnm,
			"xqm": xqm,
		}).
		Post(scheduleURL)
	if err != nil {
		return nil, fmt.Errorf("查询课表失败: %w", err)
	}

	body := resp.String()

	if strings.Contains(body, `name="mm"`) || strings.Contains(body, `login_slogin`) {
		return nil, fmt.Errorf("session 无效，请重新登录")
	}

	return parseScheduleResponse(body)
}

func parseScheduleResponse(body string) ([]ScheduleCourse, error) {
	var result struct {
		KbList      []map[string]interface{} `json:"kbList"`
		TotalResult interface{}              `json:"totalResult"`
		PageTotal   interface{}              `json:"pageTotal"`
	}
	if err := json.Unmarshal([]byte(body), &result); err == nil && len(result.KbList) > 0 {
		courses := make([]ScheduleCourse, 0, len(result.KbList))
		for _, item := range result.KbList {
			courses = append(courses, ScheduleCourse{
				KCMC: utils.Stringify(item["kcmc"]),
				JSMC: utils.Stringify(item["cdmc"]),
				JSXM: utils.Stringify(item["xm"]),
				XQ:   utils.Stringify(item["xqj"]),
				KSJC: utils.Stringify(item["jc"]),
				JSJC: "",
				ZCMC: utils.Stringify(item["zcd"]),
			})
		}
		return courses, nil
	}

	return parseScheduleHTML(body)
}

func parseScheduleHTML(html string) ([]ScheduleCourse, error) {
	var courses []ScheduleCourse

	courseRe := regexp.MustCompile(`<div[^>]*class="[^"]*kbcontent[^"]*"[^>]*>(.*?)</div>`)
	tagRe := regexp.MustCompile(`<[^>]*>`)

	for _, match := range courseRe.FindAllStringSubmatch(html, -1) {
		content := strings.TrimSpace(tagRe.ReplaceAllString(match[1], ""))
		if content == "" {
			continue
		}
		parts := strings.Split(content, "<br/>")
		course := ScheduleCourse{}
		if len(parts) > 0 {
			course.KCMC = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			course.JSMC = strings.TrimSpace(parts[1])
		}
		if len(parts) > 2 {
			course.JSXM = strings.TrimSpace(parts[2])
		}
		courses = append(courses, course)
	}

	return courses, nil
}
