package wechat

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/esap/wechat/util"
)

// AgentsMap 应用代理，主要用于企业号
var (
	FetchDelay time.Duration = 5 * time.Minute // 默认5分钟同步一次
)

// AccessToken 回复体
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	WxErr
}

// GetAccessToken 读取AccessToken
func (s *Server) GetAccessToken() string {
	if s.accessToken == nil || s.accessToken.ExpiresIn < time.Now().Unix() {
		for i := 0; i < 3; i++ {
			err := s.getAccessToken()
			if err != nil {
				log.Printf("GetAccessToken[%v] err:%v", s.AgentId, err)
				continue
			}
			break
		}
	}
	return s.accessToken.AccessToken
}

func (s *Server) getAccessToken() (err error) {
	url := fmt.Sprintf(s.TokenUrl, s.AppId, s.Secret)
	at := new(AccessToken)
	if err = util.GetJson(url, at); err != nil {
		return
	}
	if at.ErrCode > 0 {
		return errors.New(at.ErrMsg)
	}
	Printf("[%v::%v]:%+v", s.AppId, s.AgentId, *at)
	at.ExpiresIn = time.Now().Unix() + 600
	s.accessToken = at
	return
}
