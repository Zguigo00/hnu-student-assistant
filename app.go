package main

import (
	"context"
	"fmt"

	"hnu-student-assistant/backend/auth"
	"hnu-student-assistant/backend/grade"
	"hnu-student-assistant/backend/portal"
	"hnu-student-assistant/backend/sc"
	"hnu-student-assistant/backend/schedule"

	"github.com/go-resty/resty/v2"
)

type App struct {
	ctx          context.Context
	gradeService *grade.GradeService
	scheduleSvc  *schedule.ScheduleService
	portalClient *resty.Client
	scClient     *resty.Client
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Login 登录教务系统（RSA 加密）
func (a *App) Login(username, password string) map[string]interface{} {
	client, err := auth.JwxtLogin(username, password)
	if err != nil {
		return errResp(err.Error())
	}

	a.gradeService = grade.NewGradeService(client)
	a.scheduleSvc = schedule.NewScheduleService(client)

	return okResp("登录成功", nil)
}

// GetGrades 查询成绩
func (a *App) GetGrades(xnxq string) map[string]interface{} {
	if a.gradeService == nil {
		return errResp("请先登录")
	}

	grades, err := a.gradeService.GetGrades(xnxq)
	if err != nil {
		return errResp(err.Error())
	}

	gpa := grade.CalculateGPA(grades)
	return okResp("", map[string]interface{}{
		"data":  grades,
		"gpa":   fmt.Sprintf("%.2f", gpa),
		"count": len(grades),
	})
}

// GetSchedule 查询课表
func (a *App) GetSchedule(xnxq string) map[string]interface{} {
	if a.scheduleSvc == nil {
		return errResp("请先登录")
	}

	courses, err := a.scheduleSvc.GetSchedule(xnxq)
	if err != nil {
		return errResp(err.Error())
	}

	return okResp("", map[string]interface{}{
		"data":  courses,
		"count": len(courses),
	})
}

// IsLoggedIn 检查教务系统登录状态
func (a *App) IsLoggedIn() bool {
	return a.gradeService != nil
}

// Logout 退出教务系统登录
func (a *App) Logout() {
	a.gradeService = nil
	a.scheduleSvc = nil
}

// PortalLogin 登录校园信息门户（AES-CBC 加密）
func (a *App) PortalLogin(username, password string) map[string]interface{} {
	client, err := portal.PortalLogin(username, password)
	if err != nil {
		return errResp(err.Error())
	}

	a.portalClient = client
	return okResp("门户登录成功", nil)
}

// IsPortalLoggedIn 检查门户登录状态
func (a *App) IsPortalLoggedIn() bool {
	return a.portalClient != nil
}

// GetNewsByCategory 按分类获取校园新闻
func (a *App) GetNewsByCategory(category string, page, size int) map[string]interface{} {
	if a.portalClient == nil {
		return errResp("请先登录校园门户")
	}

	newsSvc := portal.NewNewsService(a.portalClient)
	items, count, err := newsSvc.GetNewsByCategory(category, page, size)
	if err != nil {
		return errResp(err.Error())
	}

	return okResp("", map[string]interface{}{
		"data":  items,
		"count": count,
	})
}

// SCLogin 登录第二课堂系统（CAS 认证，与校园门户共用账号密码）
func (a *App) SCLogin(username, password string) map[string]interface{} {
	client, err := sc.SCLogin(username, password)
	if err != nil {
		return errResp(err.Error())
	}

	a.scClient = client
	return okResp("第二课堂登录成功", nil)
}

// IsSCLoggedIn 检查第二课堂登录状态
func (a *App) IsSCLoggedIn() bool {
	return a.scClient != nil
}

// GetSecondClassScores 查询第二课堂成绩
func (a *App) GetSecondClassScores(startTime, endTime string) map[string]interface{} {
	if a.scClient == nil {
		return errResp("请先登录第二课堂系统")
	}

	scoreSvc := sc.NewScoreService(a.scClient)
	items, err := scoreSvc.GetScores(startTime, endTime)
	if err != nil {
		return errResp(err.Error())
	}

	return okResp("", map[string]interface{}{
		"data":  items,
		"count": len(items),
	})
}

func okResp(message string, extra map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{"success": true}
	if message != "" {
		resp["message"] = message
	}
	for k, v := range extra {
		resp[k] = v
	}
	return resp
}

func errResp(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
	}
}
