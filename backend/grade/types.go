package grade

// Grade 成绩信息
type Grade struct {
	XNXQ   string `json:"xnxq"`   // 学年学期
	KCMC   string `json:"kcmc"`   // 课程名称
	KCDM   string `json:"kcdm"`   // 课程代码
	XF     string `json:"xf"`     // 学分
	CJ     string `json:"cj"`     // 成绩
	CX     string `json:"cx"`     // 重修
	KCXZDM string `json:"kcxzdm"` // 课程性质代码
	KCXZMC string `json:"kcxzmc"` // 课程性质名称
	JD     string `json:"jd"`     // 绩点
	KKXY   string `json:"kkxy"`   // 开课学院
}
