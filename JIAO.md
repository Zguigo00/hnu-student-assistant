# 构建、分发与更新指南

## 一、生成安装包

### 1. 配置 wails.json

在 `wails.json` 中添加产品信息和 NSIS 配置：

```json
{
  "info": {
    "companyName": "Zguigo",
    "productName": "淮北师范大学学生助手",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026 Zguigo",
    "comments": "查询成绩、课表、校园新闻的一站式桌面工具"
  },
  "nsis": {
    "language": "SimpChinese",
    "installerIcon": "build/windows/icon.ico",
    "uninstallerIcon": "build/windows/icon.ico"
  }
}
```

每次发版前更新 `productVersion`（如 `"1.0.1"`）。

### 2. 构建命令

```bash
# 普通构建（生成单个 exe）
wails build

# 生成 NSIS 安装包（推荐）
wails build --nsis
```

产出位置：`build/bin/hnu-student-assistant-installer.exe`

### 3. 安装包效果

- 中文安装向导界面
- 自动安装 WebView2 运行时（用户电脑没有时）
- 创建桌面快捷方式 + 开始菜单
- 自带卸载程序（控制面板可见）

---

## 二、分发给用户

1. 把 `build/bin/` 下的 `-installer.exe` 发给用户（zip 打包或直接发）
2. 用户双击 → 安装 → 桌面出现快捷方式 → 完成
3. 用户数据（localStorage）保存在用户 AppData 中，覆盖安装不会丢失

### 关于签名

未签名的 exe 安装时 Windows SmartScreen 会弹"未知发布者"警告，点"仍要运行"即可。

要消除警告需要：
1. 购买代码签名证书（约几百元/年）
2. 用 `signtool` 签名：`signtool sign /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 /f cert.pfx /p password installer.exe`
3. NSIS 脚本中取消注释 `finalize` 行也可自动签名

---

## 三、版本更新流程

### 方式一：手动分发（用户少时够用）

1. 改代码，更新 `wails.json` 中的 `productVersion`
2. `wails build --nsis`
3. 把新安装包发给用户，覆盖安装即可

### 方式二：GitHub Actions 自动构建（推荐）

在项目根目录创建 `.github/workflows/release.yml`：

```yaml
name: Build and Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Build installer
        run: wails build --nsis

      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          files: build/bin/*-installer.exe
```

发版步骤：
```bash
# 更新 wails.json 中的 productVersion
git add -A
git commit -m "v1.1.0"
git tag v1.1.0
git push origin v1.1.0
# → Actions 自动构建，自动创建 Release，自动上传安装包
```

用户去 GitHub Release 页面下载最新的 installer 即可。

### 方式三：应用内检查更新（用户多时使用）

在应用中加一个"检查更新"按钮，调用 GitHub API 比对版本：

```go
// 后端添加方法
func (a *App) CheckUpdate(currentVersion string) map[string]interface{} {
    resp, _ := resty.New().R().Get("https://api.github.com/repos/你的用户名/仓库名/releases/latest")
    // 解析 tag_name 与 currentVersion 比较
    // 返回 {hasUpdate, downloadUrl, version, notes}
}
```

前端在设置页或顶栏加"检查更新"按钮，有新版时提示用户下载。

---

## 四、快速发版清单

- [ ] 更新 `wails.json` 中的 `productVersion`
- [ ] `wails build --nsis` 验证构建成功
- [ ] 双击安装包测试：安装 → 打开 → 核心功能正常 → 卸载
- [ ] 分发安装包 / push tag 触发自动构建
