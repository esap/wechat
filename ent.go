package wechat

import (
	"encoding/base64"
	"errors"
	"time"
)

// WXAPI_ENT 企业号接口
const (
	WXAPI_ENT          = "https://qyapi.weixin.qq.com/cgi-bin/"
	WXAPI_TOKEN_ENT    = WXAPI_ENT + "gettoken?corpid=%s&corpsecret=%s"
	WXAPI_MSG_ENT      = WXAPI_ENT + "message/send?access_token="
	WXAPI_IP_ENT       = WXAPI_ENT + "getcallbackip?access_token="
	WXAPI_UPLOAD_ENT   = WXAPI_ENT + "media/upload?access_token=%s&type=%s"
	WXAPI_GETMEDIA_ENT = WXAPI_ENT + "media/get?access_token=%s&media_id=%s"
)

// SetEnt 初始化企业号，设置token,corpid,secrat,aesKey
func SetEnt(tk, id, sec, key string) (err error) {
	if len(key) != 43 {
		return errors.New("非法的AesKey")
	}
	token, appId, secret, safeMode, entMode = tk, id, sec, true, true
	msgUrl = WXAPI_MSG_ENT
	uploadUrl = WXAPI_UPLOAD_ENT
	tokenUrl = WXAPI_TOKEN_ENT
	getMedia = WXAPI_GETMEDIA_ENT
	aesKey, err = base64.StdEncoding.DecodeString(key + "=")
	if err != nil {
		return
	}
	FetchUserList()
	return nil
}

// FetchUserList 定期获取AccessToken
func FetchUserList() {
	go func() {
		for {
			if SyncDeptList() == nil {
				SyncUserList()
			}
			time.Sleep(fetchDelay)
		}
	}()
}
