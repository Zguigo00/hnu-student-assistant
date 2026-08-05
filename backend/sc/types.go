package sc

// SecondClassScore 第课堂成绩
type SecondClassScore struct {
	UserCode  string  `json:"userCode"`
	UserName  string  `json:"userName"`
	Sxylyagrx float64 `json:"sxylyagrx"` // 思想政治与品德
	Xskjycxcy float64 `json:"xskjycxcy"` // 专业技能与创新创业
	Tydlyydjn float64 `json:"tydlyydjn"` // 体育健身运动
	Rwskyyssy float64 `json:"rwskyyssy"` // 文化艺术修养
	Shsjyzyfw float64 `json:"shsjyzyfw"` // 志愿服务与劳动实践
}
