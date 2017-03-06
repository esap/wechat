package wechat

import (
	"encoding/json"
	"time"

	"github.com/esap/wechat/util"
)

// SendMsg 发送消息
func SendMsg(v interface{}) *WxErr {
	time.Sleep(1 * time.Second)
	url := msgUrl + GetAccessToken()
	body, err := util.PostJson(url, v)
	if err != nil {
		return &WxErr{-1, err.Error()}
	}
	rst := new(WxErr)
	err = json.Unmarshal(body, rst)
	if err != nil {
		return &WxErr{-1, err.Error()}
	}
	return rst

}

// SendText 发送客服text消息
func SendText(to string, id int, msg ...string) *WxErr {
	return SendMsg(NewText(to, id, msg...))
}

// SendImage 发送客服Image消息
func SendImage(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewImage(to, id, mediaId))
}

// SendVoice 发送客服Voice消息
func SendVoice(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewVoice(to, id, mediaId))
}

// SendFile 发送客服File消息
func SendFile(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewFile(to, id, mediaId))
}

// SendVideo 发送客服Video消息
func SendVideo(to string, id int, mediaId, title, desc string) *WxErr {
	return SendMsg(NewVideo(to, id, mediaId, title, desc))
}

// SendMusic 发送客服Music消息
func SendMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) *WxErr {
	return SendMsg(NewMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl))
}

// SendNews 发送客服news消息
func SendNews(to string, id int, arts ...Article) *WxErr {
	return SendMsg(NewNews(to, id, arts...))
}

// SendMpNews 发送加密新闻mpnews消息(仅企业号可用)
func SendMpNews(to string, id int, arts ...MpArticle) *WxErr {
	return SendMsg(NewMpNews(to, id, arts...))
}
