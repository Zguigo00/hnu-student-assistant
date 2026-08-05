package grade

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"hnu-student-assistant/backend/utils"

	"github.com/go-resty/resty/v2"
)

const gradeURL = "https://jwglxt.chnu.edu.cn/cjcx/cjcx_cxXsgrcj.html"

// GradeService 成绩查询服务。
type GradeService struct {
	httpClient *resty.Client
}

func NewGradeService(client *resty.Client) *GradeService {
	return &GradeService{httpClient: client}
}

// GetGrades 查询成绩。xnxq 格式如 "2024-2025-1"，为空则查询全部。
func (s *GradeService) GetGrades(xnxq string) ([]Grade, error) {
	xnm, xqm := utils.ParseXnxq(xnxq)

	resp, err := s.httpClient.R().
		SetQueryParam("doType", "query").
		SetQueryParam("gnmkdm", "N305005").
		SetFormData(map[string]string{
			"xnm":                    xnm,
			"xqm":                    xqm,
			"_search":                "false",
			"nd":                     "",
			"queryModel.showCount":   "100",
			"queryModel.currentPage": "1",
			"queryModel.sortName":    "",
			"queryModel.sortOrder":   "asc",
		}).
		Post(gradeURL)
	if err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
	}

	body := resp.String()

	if strings.Contains(body, `name="mm"`) || strings.Contains(body, `login_slogin`) {
		return nil, fmt.Errorf("session 无效，请重新登录")
	}

	return parseGradeResponse(body)
}

// CalculateGPA 计算加权平均绩点。
func CalculateGPA(grades []Grade) float64 {
	totalPoints := 0.0
	totalCredits := 0.0
	for _, g := range grades {
		credits := utils.ParseFloat(g.XF)
		points := utils.ParseFloat(g.JD)
		if credits > 0 && points >= 0 {
			totalPoints += credits * points
			totalCredits += credits
		}
	}
	if totalCredits == 0 {
		return 0
	}
	return totalPoints / totalCredits
}

func parseGradeResponse(body string) ([]Grade, error) {
	var result struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return parseGradeHTML(body)
	}

	if len(result.Items) > 0 {
		grades := make([]Grade, 0, len(result.Items))
		for _, item := range result.Items {
			grades = append(grades, Grade{
				XNXQ:   utils.Stringify(item["xnmmc"]) + "-" + utils.Stringify(item["xqmmc"]),
				KCMC:   utils.Stringify(item["kcmc"]),
				KCDM:   utils.Stringify(item["kch_id"]),
				XF:     utils.Stringify(item["xf"]),
				CJ:     utils.Stringify(item["cj"]),
				CX:     utils.Stringify(item["cxbj"]),
				KCXZDM: utils.Stringify(item["kcxzdm"]),
				KCXZMC: utils.Stringify(item["kcxzmc"]),
				JD:     utils.Stringify(item["jd"]),
				KKXY:   utils.Stringify(item["kkbmmc"]),
			})
		}
		return grades, nil
	}

	return parseGradeHTML(body)
}

func parseGradeHTML(html string) ([]Grade, error) {
	var grades []Grade

	rowRe := regexp.MustCompile(`<tr[^>]*>(.*?)</tr>`)
	cellRe := regexp.MustCompile(`<td[^>]*>(.*?)</td>`)
	tagRe := regexp.MustCompile(`<[^>]*>`)

	for _, row := range rowRe.FindAllStringSubmatch(html, -1) {
		cells := cellRe.FindAllStringSubmatch(row[1], -1)
		if len(cells) >= 6 {
			clean := func(s string) string {
				return strings.TrimSpace(tagRe.ReplaceAllString(s, ""))
			}
			grades = append(grades, Grade{
				XNXQ:   clean(cells[0][1]),
				KCMC:   clean(cells[1][1]),
				XF:     clean(cells[2][1]),
				CJ:     clean(cells[3][1]),
				JD:     clean(cells[4][1]),
				KCXZMC: clean(cells[5][1]),
			})
		}
	}

	return grades, nil
}
