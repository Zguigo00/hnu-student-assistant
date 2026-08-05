# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

淮北师范大学学生助手 — 桌面应用，通过正方教务系统 (jwglxt.chnu.edu.cn) 查询成绩和课表，通过信息门户 (oshall.chnu.edu.cn) 查看校园新闻。

## Commands

```bash
wails dev          # 开发模式（热重载）
wails build        # 生产构建，输出到 build/bin/
npm run build      # 仅前端构建（vue-tsc + vite）
npm run dev        # 仅前端开发服务器
```

## Architecture

**Wails v2** = Go 后端 + Vue 3 前端。Go 方法通过 Wails 自动生成的 JS 绑定暴露给前端。

### 后端 (Go)

```
main.go                          → Wails 入口，嵌入 frontend/dist
app.go                           → App 结构体，暴露给前端的 8 个方法 + okResp/errResp 响应构建器
backend/
  auth/login.go                  → 教务系统 RSA 登录流程（URL 常量内联）
  grade/
    types.go                     → Grade 数据模型
    service.go                   → 成绩查询 + GPA 计算（URL 常量内联）
  schedule/
    types.go                     → ScheduleCourse 数据模型
    service.go                   → 课表查询（URL 常量内联）
  portal/
    types.go                     → NewsItem 数据模型
    urls.go                      → CAS 登录 URL、新闻 API URL
    login.go                     → 校园门户 CAS 登录（AES-CBC 加密）
    news.go                      → 校园新闻查询
  utils/converter.go             → ParseXnxq 学期解析、ParseFloat、Stringify、ExtractHiddenFields
```

- 所有 API 方法通过 `okResp`/`errResp` 返回统一格式：`{success, message, data}`
- 教务系统响应可能是 JSON 或 HTML，解析策略为 JSON 优先、HTML 回退
- JSON 解析用 `map[string]interface{}`（正方系统混合字符串和数字字段）
- URL 有查询参数：成绩 `?doType=query&gnmkdm=N305005`，课表 `?gnmkdm=N2151`

### 前端 (Vue 3 + TypeScript)

```
src/
  api/jwxt.ts        → 封装 Wails 教务系统登录/成绩/课表绑定
  api/portal.ts      → 封装 Wails 门户登录/新闻绑定
  types/index.ts     → 与后端 models 对应的 TS 接口
  constants/index.ts → 学期列表、周次、节次、配色池
  stores/jwxt.ts     → Pinia 教务系统状态 + 凭证管理
  stores/portal.ts   → Pinia 门户状态 + 凭证管理
  router/index.ts    → Hash 路由
  layouts/           → 侧边栏 + 顶栏布局
  pages/             → Grades、Schedule、News、Settings
```

- 前端不直接 import `wailsjs/go/main/App`，统一通过 `src/api/` 层调用
- 新增页面：在 `router/index.ts` 添加路由，在 `MainLayout.vue` 添加菜单项
- 两套账号（教务系统/校园门户）均在 Settings 页面管理，访问时通过 store.ensureLogin() 自动登录

### 教务系统登录流程（RSA）

1. GET 登录页 → 提取 csrftoken + 所有隐藏字段
2. GET `/xtgl/login_getPublicKey.html` → 解析 RSA 公钥 (modulus + exponent)
3. RSA PKCS1v15 加密密码 → Base64 编码
4. POST 登录页 → 提交 yhm、mm、csrftoken 及所有隐藏字段

### 校园门户登录流程（AES-CBC）

1. GET CAS 登录页 (`auth.chnu.edu.cn/authserver/login?service=...`) → 提取隐藏字段 (lt, execution, _eventId)
2. 从页面 `<script>` 中提取 AES 密钥参数 (pwdDefaultEncryptSalt)
3. AES-CBC 加密密码：UTF-8 → Zero Padding → AES-CBC → Base64
4. POST 登录表单 → resty 自动跟随重定向获取 session cookies
5. 门户凭证由前端 Settings 页面保存到 localStorage，访问新闻时自动登录

### 新增查询功能的模式

1. `backend/<功能名>/types.go` 添加数据模型
2. `backend/<功能名>/service.go` 新建服务文件（URL 常量内联），JSON 优先 + HTML 回退解析
3. `app.go` 添加方法暴露给前端（使用 `okResp`/`errResp`）
4. `frontend/src/api/jwxt.ts` 添加封装函数
5. `frontend/src/types/index.ts` 添加 TS 接口
6. `frontend/src/pages/` 新建页面组件

## 注意事项

- 正方教务系统 HTML 的 `value` 属性可能无引号（`value= 1` 而非 `value="1"`），正则需处理空格
- 正方 JSON 可能含数字字段（如 `pageTotal:0`），用 `interface{}` 接收而非 `string`
- 学期格式 "2024-2025-1" 转为教务系统参数：xnm="2024"，xqm="3"（第一学期）或 "12"（第二学期）
- `frontend/wailsjs/` 是 Wails 自动生成的，不要手动修改
- 教务系统和校园门户是两套独立的认证体系，账号密码不同，session 互不影响
- `extractCryptoParams` 中的 AES 参数变量名（如 `pwdDefaultEncryptSalt`）需要根据实际 CAS 登录页确认
