package wechat

import (
	"fmt"
	"net/url"

	"github.com/esap/wechat/util"
)

// OAUTH2PAGE oauth2鉴权
const (
	OAUTH2PAGE    = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=snsapi_base&state=110#wechat_redirect"
	Code2PAGE     = "https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code"
	Code2PAGE_Ent = "https://qyapi.weixin.qq.com/cgi-bin/miniprogram/jscode2session?access_token=%v&js_code=%v&grant_type=authorization_code"
)

type WxSession struct {
	WxErr
	SessionKey string `json:"session_key"`
	// corp
	CorpId string `json:"corpid"`
	UserId string `json:"userid"`
	// mp
	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}

// GetOauth2Url 获取鉴权页面
func GetOauth2Url(corpId, host string) string {
	return fmt.Sprintf(OAUTH2PAGE, corpId, url.QueryEscape(host))
}

// Jscode2Session code换session
func (s *Server) Jscode2Session(code string) (ws *WxSession, err error) {
	url := fmt.Sprintf(Code2PAGE, s.AppId, s.Secret, code)
	ws = new(WxSession)
	err = util.GetJson(url, ws)

	if ws.Error() != nil {
		err = ws.Error()
	}
	return
}

// Jscode2SessionEnt code换session（企业微信）
func (s *Server) Jscode2SessionEnt(code string) (ws *WxSession, err error) {
	url := fmt.Sprintf(Code2PAGE_Ent, s.GetAccessToken(), code)
	ws = new(WxSession)
	err = util.GetJson(url, ws)

	if ws.Error() != nil {
		err = ws.Error()
	}
	return
}
