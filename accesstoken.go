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
	if s.accessToken == nil || s.accessToken.ExpiresIn < time.Now().Unix() {
		for i := 0; i < 3; i++ {
			err := s.getAccessToken()
			if err != nil {
				log.Printf("GetAccessToken[%v] err:%v", s.AgentId, err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
	return s.accessToken.AccessToken
}

// GetUserAccessToken 获取通讯录AccessToken
func (s *Server) GetUserAccessToken() string {
	if us, ok := UserServerMap[s.AppId]; ok {
		return us.GetAccessToken()
	}
	return s.GetAccessToken()
}

func (s *Server) getAccessToken() (err error) {
	url := fmt.Sprintf(s.TokenUrl, s.AppId, s.Secret)
	at := new(AccessToken)
	if err = util.GetJson(url, at); err != nil {
		return
	}
	if at.ErrCode > 0 {
		return at.Error()
	}
	Printf("[%v::%v]:%+v", s.AppId, s.AgentId, *at)
	at.ExpiresIn = time.Now().Unix() + 500
	s.accessToken = at
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
	log.Printf("[%v::%v-JsApi] >>> %+v", s.AppId, s.AgentId, *at)
	at.ExpiresIn = time.Now().Unix() + 500
	s.ticket = at
	return
}

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

func (s *Server) GetJsConfig(Url string) *JsConfig {
	jc := &JsConfig{Beta: true, Debug: Debug, AppId: s.AppId}
	jc.Timestamp = time.Now().Unix()
	jc.Nonsestr = "esap"
	jc.Signature = sortSha1(fmt.Sprintf("jsapi_ticket=%v&noncestr=%v&timestamp=%v&url=%v", s.GetTicket(), jc.Nonsestr, jc.Timestamp, Url))
	jc.JsApiList = []string{"scanQRCode"}
	jc.Url = Url
	jc.App = s.AgentId
	Println("jsconfig:", jc) //debug
	return jc
}
