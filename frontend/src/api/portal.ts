import { PortalLogin, IsPortalLoggedIn, GetNewsByCategory } from '../../wailsjs/go/main/App'
import type { PortalLoginResult, NewsResult } from '../types'

export async function portalLogin(username: string, password: string): Promise<PortalLoginResult> {
  return (await PortalLogin(username, password)) as PortalLoginResult
}

export async function isPortalLoggedIn(): Promise<boolean> {
  return await IsPortalLoggedIn()
}

export async function getNewsByCategory(category: string, page: number, size: number): Promise<NewsResult> {
  return (await GetNewsByCategory(category, page, size)) as NewsResult
}
