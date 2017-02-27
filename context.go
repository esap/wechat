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
	Resp      interface{}
	Writer    http.ResponseWriter
	Request   *http.Request
	repCount  int
}

// Reply 被动回复消息
func (c *Context) Reply() *Context {
	if c.Request.Method != "POST" || c.repCount > 0 {
		return c
	}
	c.write(c.Resp)
	c.repCount++
	return c
}

// ReplySuccess 如果不能在5秒内处理完，应该先回复success，然后通过客服消息通知用户
func (c *Context) ReplySuccess() *Context {
	c.Writer.Write([]byte("success"))
	return c
}

// Send 主动发送消息(客服)
func (c *Context) Send() *Context {
	SendMsg(c.Resp)
	return c
}

func (c *Context) write(v interface{}) {
	b, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println("write()->xmlMarsha1 error:", err)
		c.Writer.Write([]byte{})
	}
	Println(string(b))
	if safeMode {
		b, err = EncryptMsg(b, c.Timestamp, c.Nonce)
		if err != nil {
			log.Println("write()->EncryptMsg error:", err)
			c.Writer.Write([]byte("success"))
		}
	}
	c.Writer.Header().Set("Content-Type", "text/xml")
	c.Writer.Write(b)
}

func (c *Context) newResp(msgType string) wxResp {
	return wxResp{
		FromUserName: CDATA(c.Msg.ToUserName),
		ToUserName:   CDATA(c.Msg.FromUserName),
		MsgType:      CDATA(msgType),
		CreateTime:   time.Now().Unix(),
		AgentId:      c.Msg.AgentID,
	}
}

// NewText Text消息
func (c *Context) NewText(text ...string) *Context {
	c.Resp = &Text{
		wxResp:  c.newResp("text"),
		content: content{CDATA(strings.Join(text, ""))}}
	return c
}

// NewImage Image消息
func (c *Context) NewImage(mediaId string) *Context {
	c.Resp = &Image{
		wxResp: c.newResp("image"),
		Image:  media{CDATA(mediaId)}}
	return c
}

// NewVoice Voice消息
func (c *Context) NewVoice(mediaId string) *Context {
	c.Resp = &Voice{
		wxResp: c.newResp("voice"),
		Voice:  media{CDATA(mediaId)}}
	return c
}

// NewFile File消息
func (c *Context) NewFile(mediaId string) *Context {
	c.Resp = &File{
		wxResp: c.newResp("file"),
		File:   media{CDATA(mediaId)}}
	return c
}

// NewVideo Video消息
func (c *Context) NewVideo(mediaId, title, desc string) *Context {
	c.Resp = &Video{
		wxResp: c.newResp("video"),
		Video:  video{CDATA(mediaId), CDATA(title), CDATA(desc)}}
	return c
}

// NewNews News消息
func (c *Context) NewNews(arts ...Article) *Context {
	new := News{
		wxResp:       c.newResp("news"),
		ArticleCount: len(arts),
	}
	new.Articles.Item = arts
	c.Resp = &new
	return c
}

// NewMusic Music消息
func (c *Context) NewMusic(mediaId, title, desc, musicUrl, hqMusicUrl string) *Context {
	c.Resp = &Music{
		wxResp: c.newResp("music"),
		Music:  music{CDATA(mediaId), CDATA(title), CDATA(desc), CDATA(musicUrl), CDATA(hqMusicUrl)}}
	return c
}
