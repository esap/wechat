package wechat

import (
	"encoding/json"

	"github.com/esap/wechat/util"
)

// MsgToDo 队列消息
type MsgToDo struct {
	Msg     interface{}
	AgentId int
}

// MsgQueue 主动消息队列
var MsgQueue chan *MsgToDo

// MsgQueueAdd 添加队列消息消息
func MsgQueueAdd(v interface{}, ag ...int) {
	agent := 0
	if len(ag) > 0 {
		agent = ag[0]
	}
	MsgQueue <- &MsgToDo{v, agent}
}

func init() {
	MsgQueue = make(chan *MsgToDo, 10000)
	go func() {
		for {
			msg := <-MsgQueue
			SendMsg(msg.Msg, msg.AgentId)
		}
	}()
}

// SendMsg 发送消息
func SendMsg(v interface{}, ag ...int) *WxErr {
	agent := 0
	if len(ag) > 0 {
		agent = ag[0]
	}
	at, err := GetAgentAccessToken(agent)
	if err != nil {
		return &WxErr{-1, err.Error()}
	}
	url := msgUrl + at
	body, err := util.PostJson(url, v)
	if err != nil {
		return &WxErr{-1, err.Error()}
	}
	rst := new(WxErr)
	err = json.Unmarshal(body, rst)
	if err != nil {
		return &WxErr{-1, err.Error()}
	}
	Printf("发送消息:%+v\n回执:%+v", v, *rst)
	return rst

}

// SendText 发送客服text消息
func SendText(to string, id int, msg ...string) *WxErr {
	return SendMsg(NewText(to, id, msg...), id)
}

// SendImage 发送客服Image消息
func SendImage(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewImage(to, id, mediaId), id)
}

// SendVoice 发送客服Voice消息
func SendVoice(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewVoice(to, id, mediaId), id)
}

// SendFile 发送客服File消息
func SendFile(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewFile(to, id, mediaId), id)
}

// SendVideo 发送客服Video消息
func SendVideo(to string, id int, mediaId, title, desc string) *WxErr {
	return SendMsg(NewVideo(to, id, mediaId, title, desc), id)
}

// SendTextcard 发送客服extcard消息
func SendTextcard(to string, id int, title, desc, url string) *WxErr {
	return SendMsg(NewTextcard(to, id, title, desc, url), id)
}

// SendMusic 发送客服Music消息
func SendMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) *WxErr {
	return SendMsg(NewMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl), id)
}

// SendNews 发送客服news消息
func SendNews(to string, id int, arts ...Article) *WxErr {
	return SendMsg(NewNews(to, id, arts...), id)
}

// SendMpNews 发送加密新闻mpnews消息(仅企业号可用)
func SendMpNews(to string, id int, arts ...MpArticle) *WxErr {
	return SendMsg(NewMpNews(to, id, arts...), id)
}

// SendMpNews2 发送加密新闻mpnews消息(直接使用mediaId)
func SendMpNews2(to string, id int, mediaId string) *WxErr {
	return SendMsg(NewMpNews2(to, id, mediaId), id)
}
