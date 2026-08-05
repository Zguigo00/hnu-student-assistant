// 学期列表，新增学期只需在此添加。
export const SEMESTERS = [
  '2024-2025-1',
  '2024-2025-2',
  '2025-2026-1',
  '2025-2026-2',
]

// 默认选中的学期。
export const DEFAULT_SEMESTER = '2025-2026-2'

// 一周七天。
export const WEEK_DAYS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

// 每日节次区间。
export const SECTIONS = ['1-2', '3-4', '5-6', '7-8', '9-10']

// 周次列表（1-20周）。
export const WEEKS = Array.from({ length: 20 }, (_, i) => i + 1)

// 课表卡片配色池。
export const COURSE_COLORS = [
  '#e3f2fd', '#fce4ec', '#e8f5e9', '#fff3e0', '#f3e5f5',
  '#e0f7fa', '#fff9c4', '#efebe9', '#e8eaf6', '#fbe9e7',
]
