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
	AgentsMap                    = make(map[int]string)
	accessTokenMap               = make(map[int]*accessToken)
	fetchDelay     time.Duration = 1200 * time.Second // 默认20分钟同步一次
)

// accessToken 回复体
type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	WxErr
}

// GetAccessToken 读取AccessToken
func GetAccessToken() string {
	i := 0
	for i < 3 {
		i++
		at, err := GetAgentAccessToken(999999)
		if err != nil {
			log.Println("GetAccessToken err:", err)
			continue
		}
		return at
	}
	return ""
}

// GetAgentAccessToken 读取AgentAccessToken
func GetAgentAccessToken(agentId int) (accesstoken string, err error) {
	if _, ok := accessTokenMap[agentId]; !ok {
		accessTokenMap[agentId] = new(accessToken)
	}

	if accessTokenMap[agentId].ExpiresIn < time.Now().Unix() {
		if err = fetchAccessToken(agentId); err != nil {
			return
		}
	}
	accesstoken = accessTokenMap[agentId].AccessToken
	return
}

func fetchAccessToken(agentId int) (err error) {
	if _, ok := AgentsMap[agentId]; !ok {
		AgentsMap[agentId] = secret
	}
	url := fmt.Sprintf(tokenUrl, appId, AgentsMap[agentId])
	at := new(accessToken)
	if err = util.GetJson(url, at); err != nil {
		return
	}
	if at.ErrCode > 0 {
		return errors.New(at.ErrMsg)
	}
	Printf("AccessToken[%v]:%+v", agentId, at)
	at.ExpiresIn += time.Now().Unix() - 1
	accessTokenMap[agentId] = at
	return nil
}
