package portal

import (
	"hnu-student-assistant/backend/utils"

	"github.com/go-resty/resty/v2"
)

// PortalLogin 登录校园信息门户（CAS 认证 + AES-CBC 加密）。
// 返回已认证的 resty.Client，自动设置 X-Access-Token。
func PortalLogin(username, password string) (*resty.Client, error) {
	return utils.CASLogin(
		CasLoginURL,
		username, password,
		utils.ValidateCASTicket(
			"https://oshall.chnu.edu.cn/apiZhxy/cas/client/validateLogin",
			"https://oshall.chnu.edu.cn/zhxy-new",
		),
	)
}
