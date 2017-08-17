package wechat

import (
	"fmt"
	"net/url"
)

// OAUTH2PAGE oauth2鉴权
const (
	OAUTH2PAGE = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=snsapi_base&state=110#wechat_redirect"
)

// GetOauth2Url 获取鉴权页面
func GetOauth2Url(corpId, host string) string {
	return fmt.Sprintf(OAUTH2PAGE, corpId, url.QueryEscape(host))
}
