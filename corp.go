package wechat

import (
	"encoding/base64"
	"time"
)

// WXAPI_ENT 企业号接口
const (
	WXAPI_ENT       = "https://qyapi.weixin.qq.com/cgi-bin/"
	WXAPI_TOKEN_ENT = WXAPI_ENT + "gettoken?corpid=%s&corpsecret=%s"
	WXAPI_MSG_ENT   = WXAPI_ENT + "message/send?access_token="
)

// SetEnt 初始化企业号，设置token,corpid,secrat,aesKey
func (s *Server) SetEnt(token, appId, secret, aeskey string, agentId ...int) (err error) {
	s.Token, s.AppId, s.Secret, s.SafeMode, s.EntMode = token, appId, secret, true, true
	if len(agentId) > 0 {
		s.AgentId = agentId[0]
	}
	s.RootUrl = WXAPI_ENT
	s.MsgUrl = WXAPI_MSG_ENT
	s.TokenUrl = WXAPI_TOKEN_ENT
	if aeskey != "" {
		s.AesKey, err = base64.StdEncoding.DecodeString(aeskey + "=")
		if err != nil {
			return
		}
	}
	s.FetchUserList()
	return nil
}

// SetEnt 初始化企业号，设置token,corpid,secrat,aesKey
func SetEnt(token, appId, secret, aeskey string, agentId ...int) (err error) {
	return std.SetEnt(token, appId, secret, aeskey, agentId...)
}

// FetchUserList 定期获取AccessToken
func (s *Server) FetchUserList() {
	i := 0
	go func() {
		for {
			if s.SyncDeptList() == nil {
				if s.SyncUserList() != nil && i < 3 {
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
