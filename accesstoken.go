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
	atlock      sync.Mutex
	fetchDelay  time.Duration = 5400 * time.Second // 默认1.5小时获取一次
)

// accessToken 回复体
type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	wxErr
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
