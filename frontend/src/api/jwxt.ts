import { Login, IsLoggedIn, GetGrades, GetSchedule } from '../../wailsjs/go/main/App'
import type { LoginResult, GradeResult, ScheduleResult } from '../types'

export async function jwxtLogin(username: string, password: string): Promise<LoginResult> {
  return (await Login(username, password)) as LoginResult
}

export async function isLoggedIn(): Promise<boolean> {
  return await IsLoggedIn()
}

// 查询成绩。xnxq 格式如 "2024-2025-1"，为空则查询全部。
export async function getGrades(xnxq: string = ''): Promise<GradeResult> {
  return (await GetGrades(xnxq)) as GradeResult
}

// 查询课表。xnxq 格式如 "2024-2025-1"。
export async function getSchedule(xnxq: string): Promise<ScheduleResult> {
  return (await GetSchedule(xnxq)) as ScheduleResult
}
