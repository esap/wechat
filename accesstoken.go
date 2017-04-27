package wechat

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/esap/wechat/util"
)

var (
	accesstoken string
	tokenSvr    string
	atlock      sync.Mutex
	fetchDelay  time.Duration = 5400 * time.Second // 默认1.5小时获取一次
)

// accessToken 回复体
type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	WxErr
}

// GetAccessToken 读取AccessToken
func GetAccessToken() string {
	atlock.Lock()
	defer atlock.Unlock()
	return accesstoken
}

// FetchAccessToken 定期获取AccessToken
func FetchAccessToken(url string) {
	go func() {
		for {
			if err := fetchAccessToken(url, appId, secret); err != nil {
				log.Println("FetchAccessToken...", err)
			} else {
				if err = UpdateDeptList(); err == nil {
					err = UpdateUserList()
				}
			}
			time.Sleep(fetchDelay)
		}
	}()
}

func fetchAccessToken(url, appId, secret string) error {
	url = fmt.Sprintf(url, appId, secret)
	at := new(accessToken)
	if err := util.GetJson(url, at); err != nil {
		return err
	}
	if at.ErrCode > 0 {
		return errors.New(at.ErrMsg)
	}
	atlock.Lock()
	accesstoken = at.AccessToken
	atlock.Unlock()
	Println("AccessToken:", GetAccessToken(), fetchDelay)
	return nil
}

// SetAccessTokenSvr 设置token中心服务器
func SetAccessTokenSvr(url string) {
	tokenSvr = url
}

// GetAccessTokenSvr 获取token中心服务器
func GetAccessTokenSvr() string {
	if tokenSvr != "" {
		return tokenSvr
	}
	return WXAPI_TOKEN_ENT
}
