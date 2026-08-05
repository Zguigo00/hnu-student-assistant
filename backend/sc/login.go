package sc

import (
	"hnu-student-assistant/backend/utils"

	"github.com/go-resty/resty/v2"
)

const (
	casLoginURL      = "https://auth.chnu.edu.cn/authserver/login?service=https%3A%2F%2Fsc.chnu.edu.cn%2FhbsdWeb-prod"
	tokenValidateURL = "https://sc.chnu.edu.cn/hbsdDkApi/cas/client/checkSsoLogin"
	casService       = "https://sc.chnu.edu.cn/hbsdWeb-prod"
)

// SCLogin 登录第二课堂系统（CAS 认证，与校园门户共用账号密码）。
// 返回已设置 X-Access-Token 的 resty.Client。
func SCLogin(username, password string) (*resty.Client, error) {
	return utils.CASLogin(
		casLoginURL,
		username, password,
		utils.ValidateCASTicket(tokenValidateURL, casService),
	)
}
