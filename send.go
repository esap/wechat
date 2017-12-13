package wechat

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

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

// SendText 发送客服text消息,过长时自动拆分
func (s *Server) SendText(to string, agentId int, msg string, safe ...int) (e *WxErr) {
	if len(safe) > 0 && safe[0] == 1 {
		s.SafeOpen()
		defer s.SafeClose()
	}
	//	m := strings.Join(msg, "")
	leng := utf8.RuneCountInString(msg)
	n := leng/500 + 1

	if n == 1 {
		return s.SendMsg(s.NewText(to, agentId, msg))
	} else {
		for i := 0; i < n; i++ {
			e = s.SendMsg(s.NewText(to, agentId, fmt.Sprintf("%s\n(%v/%v)", Substr(msg, i*500, (i+1)*500), i+1, n)))
		}
	}
	return
}

// Substr 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length || end < 0 {
		return ""
	}

	if end > length {
		return string(rs[start:])
	}
	return string(rs[start:end])
}

// SendImage 发送客服Image消息
func (s *Server) SendImage(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(s.NewImage(to, id, mediaId))
}

// SendVoice 发送客服Voice消息
func (s *Server) SendVoice(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(s.NewVoice(to, id, mediaId))
}

// SendFile 发送客服File消息
func (s *Server) SendFile(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(s.NewFile(to, id, mediaId))
}

// SendVideo 发送客服Video消息
func (s *Server) SendVideo(to string, id int, mediaId, title, desc string) *WxErr {
	return s.SendMsg(s.NewVideo(to, id, mediaId, title, desc))
}

// SendTextcard 发送客服extcard消息
func (s *Server) SendTextcard(to string, id int, title, desc, url string) *WxErr {
	return s.SendMsg(s.NewTextcard(to, id, title, desc, url))
}

// SendMusic 发送客服Music消息
func (s *Server) SendMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) *WxErr {
	return s.SendMsg(s.NewMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl))
}

// SendNews 发送客服news消息
func (s *Server) SendNews(to string, id int, arts ...Article) *WxErr {
	return s.SendMsg(s.NewNews(to, id, arts...))
}

// SendMpNews 发送加密新闻mpnews消息(仅企业号可用)
func (s *Server) SendMpNews(to string, id int, arts ...MpArticle) *WxErr {
	return s.SendMsg(s.NewMpNews(to, id, arts...))
}

// SendMpNewsId 发送加密新闻mpnews消息(直接使用mediaId)
func (s *Server) SendMpNewsId(to string, id int, mediaId string) *WxErr {
	return s.SendMsg(s.NewMpNewsId(to, id, mediaId))
}

// SendMsg 发送消息
func SendMsg(v interface{}) *WxErr {
	return std.SendMsg(v)
}

// SendText 发送客服text消息
func SendText(to string, id int, msg string, safe ...int) *WxErr {
	return std.SendText(to, id, msg, safe...)
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
