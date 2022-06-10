package wechat

import (
	"fmt"
	"log"
	"time"

	"github.com/esap/wechat/util"
)

// FetchDelay 默认5分钟同步一次
var FetchDelay time.Duration = 5 * time.Minute

// AccessToken 回复体
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	WxErr
}

// GetAccessToken 读取AccessToken
func (s *Server) GetAccessToken() string {
	s.Lock()
	defer s.Unlock()
	var err error
	if s.accessToken == nil || s.accessToken.ExpiresIn < time.Now().Unix() {
		for i := 0; i < 3; i++ {
			err = s.getAccessToken()
			if err == nil {
				break
			}
			log.Printf("GetAccessToken[%v] %v", s.AgentId, err)
			time.Sleep(time.Second)
		}
		if err != nil {
			return ""
		}
	}
	return s.accessToken.AccessToken
}

// GetUserAccessToken 获取企业微信通讯录AccessToken
func (s *Server) GetUserAccessToken() string {
	if us, ok := UserServerMap[s.AppId]; ok {
		return us.GetAccessToken()
	}
	return s.GetAccessToken()
}

func (s *Server) getAccessToken() (err error) {
	if s.ExternalTokenHandler != nil {
		s.accessToken = s.ExternalTokenHandler(s.AppId, s.AppName)
		Printf("***%v[%v]远程获取token:%v", util.Substr(s.AppId, 14, 30), s.AgentId, s.accessToken)
		return
	}
	url := fmt.Sprintf(s.TokenUrl, s.AppId, s.Secret)
	at := new(AccessToken)
	if err = util.GetJson(url, at); err != nil {
		return
	}
	if at.ErrCode > 0 {
		return at.Error()
	}
	at.ExpiresIn = time.Now().Unix() + at.ExpiresIn - 5
	s.accessToken = at
	Printf("***%v[%v]本地获取token:%v", util.Substr(s.AppId, 14, 30), s.AgentId, s.accessToken)
	return

}

// Ticket JS-SDK
type Ticket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
	WxErr
}

// GetTicket 读取获取Ticket
func (s *Server) GetTicket() string {
	if s.ticket == nil || s.ticket.ExpiresIn < time.Now().Unix() {
		for i := 0; i < 3; i++ {
			err := s.getTicket()
			if err != nil {
				log.Printf("getTicket[%v] err:%v", s.AgentId, err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
	return s.ticket.Ticket
}

func (s *Server) getTicket() (err error) {
	url := s.JsApi + s.GetAccessToken()
	at := new(Ticket)
	if err = util.GetJson(url, at); err != nil {
		return
	}
	if at.ErrCode > 0 {
		return at.Error()
	}
	Printf("[%v::%v-JsApi] >>> %+v", s.AppId, s.AgentId, *at)
	at.ExpiresIn = time.Now().Unix() + 500
	s.ticket = at
	return
}

// JsConfig Jssdk配置
type JsConfig struct {
	Beta      bool     `json:"beta"`
	Debug     bool     `json:"debug"`
	AppId     string   `json:"appId"`
	Timestamp int64    `json:"timestamp"`
	Nonsestr  string   `json:"nonceStr"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
	Url       string   `json:"jsurl"`
	App       int      `json:"jsapp"`
}

// GetJsConfig 获取Jssdk配置
func (s *Server) GetJsConfig(Url string) *JsConfig {
	jc := &JsConfig{Beta: true, Debug: Debug, AppId: s.AppId}
	jc.Timestamp = time.Now().Unix()
	jc.Nonsestr = "esap"
	jc.Signature = util.SortSha1(fmt.Sprintf("jsapi_ticket=%v&noncestr=%v&timestamp=%v&url=%v", s.GetTicket(), jc.Nonsestr, jc.Timestamp, Url))
	// TODO：可加入其他apilist
	jc.JsApiList = []string{"scanQRCode"}
	jc.Url = Url
	jc.App = s.AgentId
	Println("jsconfig:", jc) // Debug
	return jc
}
