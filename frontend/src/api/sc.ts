import { SCLogin, IsSCLoggedIn, GetSecondClassScores } from '../../wailsjs/go/main/App'
import type { LoginResult, SecondClassResult } from '../types'

export async function scLogin(username: string, password: string): Promise<LoginResult> {
  return (await SCLogin(username, password)) as LoginResult
}

export async function isSCLoggedIn(): Promise<boolean> {
  return await IsSCLoggedIn()
}

export async function getSecondClassScores(startTime: string, endTime: string): Promise<SecondClassResult> {
  return (await GetSecondClassScores(startTime, endTime)) as SecondClassResult
}
