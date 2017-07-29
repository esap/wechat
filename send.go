package wechat

import (
	"encoding/json"

	"github.com/esap/wechat/util"
)

// MsgQueueAdd 添加队列消息消息
func (s *Server) MsgQueueAdd(v interface{}) {
	s.MsgQueue <- v
}

func (s *Server) init() {
	s.MsgQueue = make(chan interface{}, 10000)
	go func() {
		for {
			msg := <-s.MsgQueue
			e := s.SendMsg(msg)
			if e.ErrCode != 0 {
				Println("MsgQueueSend err:", e.ErrMsg)
			}
		}
	}()
}

// SendMsg 发送消息
func (s *Server) SendMsg(v interface{}) *WxErr {
	url := s.MsgUrl + s.GetAccessToken()
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
func (s *Server) SendText(to string, id int, msg ...string) *WxErr {
	return s.SendMsg(NewText(to, id, msg...))
}

// SendImage 发送客服Image消息
func (s *Server) SendImage(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(NewImage(to, id, mediaId))
}

// SendVoice 发送客服Voice消息
func (s *Server) SendVoice(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(NewVoice(to, id, mediaId))
}

// SendFile 发送客服File消息
func (s *Server) SendFile(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(NewFile(to, id, mediaId))
}

// SendVideo 发送客服Video消息
func (s *Server) SendVideo(to string, id int, mediaId, title, desc string) *WxErr {
	return s.SendMsg(NewVideo(to, id, mediaId, title, desc))
}

// SendTextcard 发送客服extcard消息
func (s *Server) SendTextcard(to string, id int, title, desc, url string) *WxErr {
	return s.SendMsg(NewTextcard(to, id, title, desc, url))
}

// SendMusic 发送客服Music消息
func (s *Server) SendMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) *WxErr {
	return s.SendMsg(NewMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl))
}

// SendNews 发送客服news消息
func (s *Server) SendNews(to string, id int, arts ...Article) *WxErr {
	return s.SendMsg(NewNews(to, id, arts...))
}

// SendMpNews 发送加密新闻mpnews消息(仅企业号可用)
func (s *Server) SendMpNews(to string, id int, arts ...MpArticle) *WxErr {
	return s.SendMsg(NewMpNews(to, id, arts...))
}

// SendMpNewsId 发送加密新闻mpnews消息(直接使用mediaId)
func (s *Server) SendMpNewsId(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(NewMpNewsId(to, id, mediaId))
}

// SendMsg 发送消息
func SendMsg(v interface{}) *WxErr {
	return std.SendMsg(v)
}

// SendText 发送客服text消息
func SendText(to string, id int, msg ...string) *WxErr {
	return std.SendText(to, id, msg...)
}

// SendImage 发送客服Image消息
func SendImage(to string, id int, mediaId string) *WxErr {
	return std.SendImage(to, id, mediaId)
}

// SendVoice 发送客服Voice消息
func SendVoice(to string, id int, mediaId string) *WxErr {
	return std.SendVoice(to, id, mediaId)
}

// SendFile 发送客服File消息
func SendFile(to string, id int, mediaId string) *WxErr {
	return std.SendFile(to, id, mediaId)
}

// SendVideo 发送客服Video消息
func SendVideo(to string, id int, mediaId, title, desc string) *WxErr {
	return std.SendVideo(to, id, mediaId, title, desc)
}

// SendTextcard 发送客服extcard消息
func SendTextcard(to string, id int, title, desc, url string) *WxErr {
	return std.SendTextcard(to, id, title, desc, url)
}

// SendMusic 发送客服Music消息
func SendMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) *WxErr {
	return std.SendMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl)
}

// SendNews 发送客服news消息
func SendNews(to string, id int, arts ...Article) *WxErr {
	return std.SendNews(to, id, arts...)
}

// SendMpNews 发送加密新闻mpnews消息(仅企业号可用)
func SendMpNews(to string, id int, arts ...MpArticle) *WxErr {
	return std.SendMpNews(to, id, arts...)
}

// SendMpNews2 发送加密新闻mpnews消息(直接使用mediaId)
func SendMpNewsId(to string, id int, mediaId string) *WxErr {
	return std.SendMpNewsId(to, id, mediaId)
}
