// ==================== 成绩 ====================

export interface Grade {
  xnxq: string     // 学年学期
  kcmc: string     // 课程名称
  kcdm: string     // 课程代码
  xf: string       // 学分
  cj: string       // 成绩
  cx: string       // 重修标记
  kcxzdm: string   // 课程性质代码
  kcxzmc: string   // 课程性质名称
  jd: string       // 绩点
  kkxy: string     // 开课学院
}

export interface GradeResult {
  success: boolean
  data: Grade[]
  gpa: string
  count: number
  message?: string
}

// ==================== 课表 ====================

export interface Course {
  kcmc: string     // 课程名称
  jsmc: string     // 教室名称
  jsxm: string     // 教师姓名
  xq: string       // 星期几 (1-7)
  ksjc: string     // 开始节次
  jsjc: string     // 结束节次
  zcmc: string     // 周次名称
}

export interface ScheduleResult {
  success: boolean
  data: Course[]
  count: number
  message?: string
}

// ==================== 登录 ====================

export interface LoginResult {
  success: boolean
  message: string
}

// ==================== 校园新闻 ====================

export interface NewsItem {
  id: string
  title: string
  createTime: string
  createBy: string
  look: number
  content: string
  isLinks: boolean
  type: number
}

export interface NewsResult {
  success: boolean
  data: NewsItem[]
  count: number
  message?: string
}

// ==================== 门户登录 ====================

export interface PortalLoginResult {
  success: boolean
  message: string
}

// ==================== 第课堂成绩 ====================

export interface SecondClassScore {
  userCode: string
  userName: string
  sxylyagrx: number   // 思想政治与品德
  xskjycxcy: number   // 专业技能与创新创业
  tydlyydjn: number   // 体育健身运动
  rwskyyssy: number   // 文化艺术修养
  shsjyzyfw: number   // 志愿服务与劳动实践
}

export interface SecondClassResult {
  success: boolean
  data: SecondClassScore[]
  count: number
  message?: string
}
