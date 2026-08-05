package schedule

// ScheduleCourse 课表课程
type ScheduleCourse struct {
	KCMC string `json:"kcmc"` // 课程名称
	JSMC string `json:"jsmc"` // 教室名称
	JSXM string `json:"jsxm"` // 教师姓名
	XQ   string `json:"xq"`   // 星期几 (1-7)
	KSJC string `json:"ksjc"` // 开始节次
	JSJC string `json:"jsjc"` // 结束节次
	ZCMC string `json:"zcmc"` // 周次名称
}
