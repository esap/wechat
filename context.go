package wechat

import (
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
	"time"
)

// Context 消息上下文
type Context struct {
	*Server
	Timestamp string
	Nonce     string
	Msg       *WxMsg
	Resp      interface{}
	Writer    http.ResponseWriter
	Request   *http.Request
	hasReply  bool
}

// Reply 被动回复消息
func (c *Context) Reply() (err error) {
	if c.hasReply {
		return errors.New("重复调用错误")
	}

	c.hasReply = true

	if c.Resp == nil {
		return nil
	}

	Printf("Wechat <== %+v", c.Resp)
	if c.SafeMode {
		b, err := xml.MarshalIndent(c.Resp, "", "  ")
		if err != nil {
			return err
		}
		c.Resp, err = c.EncryptMsg(b, c.Timestamp, c.Nonce)
		if err != nil {
			return err
		}
	}
	c.Writer.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	return xml.NewEncoder(c.Writer).Encode(c.Resp)
}

// Send 主动发送消息(客服)
func (c *Context) Send() *Context {
	c.AddMsg(c.Resp)
	return c
}

func (c *Context) newResp(msgType string) wxResp {
	return wxResp{
		FromUserName: CDATA(c.Msg.ToUserName),
		ToUserName:   CDATA(c.Msg.FromUserName),
		MsgType:      CDATA(msgType),
		CreateTime:   time.Now().Unix(),
		AgentId:      c.Msg.AgentID,
		Safe:         c.Safe,
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

// NewTextcard Textcard消息
func (c *Context) NewTextcard(title, description, url string) *Context {
	c.Resp = &Textcard{
		wxResp:   c.newResp(TypeTextcard),
		Textcard: textcard{CDATA(title), CDATA(description), CDATA(url)}}
	return c
}

// NewNews News消息
func (c *Context) NewNews(arts ...Article) *Context {
	news := News{
		wxResp:       c.newResp(TypeNews),
		ArticleCount: len(arts),
	}
	news.Articles.Item = arts
	c.Resp = &news
	return c
}

// NewMpNews News消息
func (c *Context) NewMpNews(mediaId string) *Context {
	news := MpNewsId{
		wxResp: c.newResp(TypeMpNews),
	}
	news.MpNews.MediaId = CDATA(mediaId)
	c.Resp = &news
	return c
}

// NewMusic Music消息
func (c *Context) NewMusic(mediaId, title, desc, musicUrl, hqMusicUrl string) *Context {
	c.Resp = &Music{
		wxResp: c.newResp(TypeMusic),
		Music:  music{CDATA(mediaId), CDATA(title), CDATA(desc), CDATA(musicUrl), CDATA(hqMusicUrl)}}
	return c
}

// Id 返回消息的来源与去向，可作为多应用管理时的用户组Id
func (c *Context) Id() string {
	return c.Msg.FromUserName + "|" + c.Msg.ToUserName
}
