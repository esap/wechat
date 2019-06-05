package wechat

import (
	"encoding/base64"
	"time"
)

// CorpAPI 企业微信接口，相关接口常量统一以此开头
const (
	CorpAPI      = "https://qyapi.weixin.qq.com/cgi-bin/"
	CorpAPIToken = CorpAPI + "gettoken?corpid=%s&corpsecret=%s"
	CorpAPIMsg   = CorpAPI + "message/send?access_token="
	CorpAPIJsapi = CorpAPI + "get_jsapi_ticket?access_token="
)

// SetEnt 初始化企业微信应用，设置token,corpid,secrat,aesKey
func (s *Server) SetEnt(token, appId, secret, aeskey string, agentId ...int) (err error) {
	s.Token, s.AppId, s.Secret, s.SafeMode, s.EntMode = token, appId, secret, true, true
	if len(agentId) > 0 {
		s.AgentId = agentId[0]
	}
	s.RootUrl = CorpAPI
	s.MsgUrl = CorpAPIMsg
	s.TokenUrl = CorpAPIToken
	s.JsApi = CorpAPIJsapi
	if aeskey != "" {
		s.AesKey, err = base64.StdEncoding.DecodeString(aeskey + "=")
		if err != nil {
			return
		}
	}
	s.FetchUserList()
	return nil
}

// FetchUserList 定期获取AccessToken
func (s *Server) FetchUserList() {
	i := 0
	go func() {
		for {
			if s.SyncDeptList() == nil {
				if s.SyncUserList() != nil && i < 2 {
					i++
					Println("尝试再次获取用户列表(", i, ")")
					continue
				}
				i = 0
			}
			s.SyncTagList()
			time.Sleep(FetchDelay)
		}
	}()
}
