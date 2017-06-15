package wechat

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"time"
)

// Context 消息上下文
type Context struct {
	Timestamp string
	Nonce     string
	Msg       *WxMsg
	MsgEnc    *WxMsgEnc
	Resp      interface{}
	Writer    http.ResponseWriter
	Request   *http.Request
	repCount  int
}

// Reply 被动回复消息
func (c *Context) Reply() *Context {
	if c.Request.Method != "POST" || c.repCount > 0 {
		log.Println("not reply...")
		return c
	}
	if safeMode {
		b, err := xml.MarshalIndent(c.Resp, "", "  ")
		if err != nil {
			log.Println("Reply()->MarshalIndent err:", err)
			c.Writer.Write([]byte{})
		}
		c.Resp, err = EncryptMsg(b, c.Timestamp, c.Nonce)
		if err != nil {
			log.Println("Reply()->EncryptMsg err:", err)
			c.Writer.Write([]byte{})
		}
	}
	Printf("reply msg:%+v", c.Resp)
	c.Writer.Header().Set("Content-Type", "text/xml")
	err := xml.NewEncoder(c.Writer).Encode(c.Resp)
	if err != nil {
		Println("Reply()->Encode err:", err)
	}
	c.repCount++
	return c
}

// ReplySuccess 如果不能在5秒内处理完，应该先回复success，然后通过客服消息通知用户
func (c *Context) ReplySuccess() *Context {
	if c.Request.Method != "POST" || c.repCount > 0 {
		return c
	}
	c.Writer.Write([]byte("success"))
	return c
}

// Send 主动发送消息(客服)
func (c *Context) Send() *Context {
	go SendMsg(c.Resp, c.Msg.AgentID)
	return c
}

func (c *Context) newResp(msgType string) wxResp {
	return wxResp{
		FromUserName: CDATA(c.Msg.ToUserName),
		ToUserName:   CDATA(c.Msg.FromUserName),
		MsgType:      CDATA(msgType),
		CreateTime:   time.Now().Unix(),
		AgentId:      c.Msg.AgentID,
		Safe:         safe,
	}
}

// NewText Text消息
func (c *Context) NewText(text ...string) *Context {
	c.Resp = &Text{
		wxResp:  c.newResp(TypeText),
		content: content{CDATA(strings.Join(text, ""))}}
	return c
}

// NewImage Image消息
func (c *Context) NewImage(mediaId string) *Context {
	c.Resp = &Image{
		wxResp: c.newResp(TypeImage),
		Image:  media{CDATA(mediaId)}}
	return c
}

// NewVoice Voice消息
func (c *Context) NewVoice(mediaId string) *Context {
	c.Resp = &Voice{
		wxResp: c.newResp(TypeVoice),
		Voice:  media{CDATA(mediaId)}}
	return c
}

// NewFile File消息
func (c *Context) NewFile(mediaId string) *Context {
	c.Resp = &File{
		wxResp: c.newResp(TypeFile),
		File:   media{CDATA(mediaId)}}
	return c
}

// NewVideo Video消息
func (c *Context) NewVideo(mediaId, title, desc string) *Context {
	c.Resp = &Video{
		wxResp: c.newResp(TypeVideo),
		Video:  video{CDATA(mediaId), CDATA(title), CDATA(desc)}}
	return c
}

// NewVideo Video消息
func (c *Context) NewTextcard(title, description, url string) *Context {
	c.Resp = &Textcard{
		wxResp:   c.newResp(TypeTextcard),
		Textcard: textcard{CDATA(title), CDATA(description), CDATA(url)}}
	return c
}

// NewNews News消息
func (c *Context) NewNews(arts ...Article) *Context {
	new := News{
		wxResp:       c.newResp(TypeNews),
		ArticleCount: len(arts),
	}
	new.Articles.Item = arts
	c.Resp = &new
	return c
}

// NewMusic Music消息
func (c *Context) NewMusic(mediaId, title, desc, musicUrl, hqMusicUrl string) *Context {
	c.Resp = &Music{
		wxResp: c.newResp(TypeMusic),
		Music:  music{CDATA(mediaId), CDATA(title), CDATA(desc), CDATA(musicUrl), CDATA(hqMusicUrl)}}
	return c
}
