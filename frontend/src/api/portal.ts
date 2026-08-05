import { PortalLogin, IsPortalLoggedIn, GetNews } from '../../wailsjs/go/main/App'
import type { PortalLoginResult, NewsResult } from '../types'

export async function portalLogin(username: string, password: string): Promise<PortalLoginResult> {
  return (await PortalLogin(username, password)) as PortalLoginResult
}

export async function isPortalLoggedIn(): Promise<boolean> {
  return await IsPortalLoggedIn()
}

export async function getNews(startIndex: number, endIndex: number): Promise<NewsResult> {
  return (await GetNews(startIndex, endIndex)) as NewsResult
}
