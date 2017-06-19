package wechat

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
)

// Type io类型汇总
const (
	TypeText     = "text"
	TypeImage    = "image"
	TypeVoice    = "voice"
	TypeMusic    = "music"
	TypeVideo    = "video"
	TypeTextcard = "textcard" // 仅企业号可用
	TypeFile     = "file"     // 仅企业号可用
	TypeNews     = "news"
	TypeMpNews   = "mpnews" // 仅企业号可用
	TypeThumb    = "thumb"
)

//WxErr 通用错误
type WxErr struct {
	ErrCode int
	ErrMsg  string
}

func (w *WxErr) Error() error {
	if w.ErrCode != 0 {
		return fmt.Errorf("err: errcode=%v , errmsg=%v", w.ErrCode, w.ErrMsg)
	}
	return nil
}

// CDATA 标准规范，XML编码成 `<![CDATA[消息内容]]>`
type CDATA string

// MarshalXML 自定义xml编码接口，实现讨论: http://stackoverflow.com/q/41951345/7493327
func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

var (
	safe int
	mu   sync.Mutex
)

//SafeOpen 打开保密消息
func SafeOpen() {
	mu.Lock()
	defer mu.Unlock()
	safe = 1
}

//SafeClose 关闭保密消息
func SafeClose() {
	mu.Lock()
	defer mu.Unlock()
	safe = 0
}

// wxResp 响应消息共用字段
// 响应消息被动回复为XML结构，文本类型采用CDATA编码规范
// 响应消息主动发送为json结构，即客服消息
type wxResp struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   CDATA    `json:"touser"`
	FromUserName CDATA    `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      CDATA    `json:"msgtype"`
	AgentId      int      `xml:"-" json:"agentid"`
	Safe         int      `xml:"-" json:"safe"`
}

func newWxResp(msgType, toUser string, agentId int) wxResp {
	return wxResp{ToUserName: CDATA(toUser), MsgType: CDATA(msgType), AgentId: agentId, Safe: safe}
}

// Text 文本消息
type Text struct {
	wxResp
	content `xml:"Content" json:"text"`
}

type content struct {
	Content CDATA `json:"content"`
}

// NewText Text 文本消息
func NewText(to string, id int, msg ...string) Text {
	return Text{
		newWxResp(TypeText, to, id),
		content{CDATA(strings.Join(msg, ""))},
	}
}

// Image 图片消息
type Image struct {
	wxResp
	Image media `json:"image"`
}

type media struct {
	MediaId CDATA `json:"media_id"`
}

// NewImage Image 消息
func NewImage(to string, id int, mediaId string) Image {
	return Image{
		newWxResp(TypeImage, to, id),
		media{CDATA(mediaId)},
	}
}

// Voice 语音消息
type Voice struct {
	wxResp
	Voice media `json:"voice"`
}

// NewVoice Voice消息
func NewVoice(to string, id int, mediaId string) Voice {
	return Voice{
		newWxResp(TypeVoice, to, id),
		media{CDATA(mediaId)},
	}
}

// File 文件消息，仅企业号支持
type File struct {
	wxResp
	File media `json:"file"`
}

// NewFile File消息
func NewFile(to string, id int, mediaId string) File {
	return File{
		newWxResp(TypeFile, to, id),
		media{CDATA(mediaId)},
	}
}

// Video 视频消息
type Video struct {
	wxResp
	Video video `json:"video"`
}

type video struct {
	MediaId     CDATA `json:"media_id"`
	Title       CDATA `json:"title"`
	Description CDATA `json:"description"`
}

// NewVideo Video消息
func NewVideo(to string, id int, mediaId, title, desc string) Video {
	return Video{
		newWxResp(TypeVideo, to, id),
		video{CDATA(mediaId), CDATA(title), CDATA(desc)},
	}
}

// Textcard 卡片消息，仅企业微信客户端有效
type Textcard struct {
	wxResp
	Textcard textcard `json:"textcard"`
}

type textcard struct {
	Title       CDATA `json:"title"`
	Description CDATA `json:"description"`
	Url         CDATA `json:"url"`
}

// NewTextcard Textcard消息
func NewTextcard(to string, id int, title, description, url string) Textcard {
	return Textcard{
		newWxResp(TypeTextcard, to, id),
		textcard{CDATA(title), CDATA(description), CDATA(url)},
	}
}

// Music 音乐消息，企业号不支持
type Music struct {
	wxResp
	Music music `json:"music"`
}

type music struct {
	Title        CDATA `json:"title"`
	Description  CDATA `json:"description"`
	MusicUrl     CDATA `json:"musicurl"`
	HQMusicUrl   CDATA `json:"hqmusicurl"`
	ThumbMediaId CDATA `json:"thumb_media_id"`
}

// NewMusic Music消息
func NewMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) Music {
	return Music{
		newWxResp(TypeMusic, to, id),
		music{CDATA(title), CDATA(desc), CDATA(musicUrl), CDATA(qhMusicUrl), CDATA(mediaId)},
	}
}

// News 新闻消息
type News struct {
	wxResp
	ArticleCount int
	Articles     struct {
		Item []Article `xml:"item" json:"articles"`
	} `json:"news"`
}

// NewNews news消息
func NewNews(to string, id int, arts ...Article) News {
	news := News{wxResp: newWxResp(TypeNews, to, id), ArticleCount: len(arts)}
	news.Articles.Item = arts
	return news
}

// Article 文章
type Article struct {
	Title       CDATA `json:"title"`
	Description CDATA `json:"description"`
	PicUrl      CDATA `json:"picurl"`
	Url         CDATA `json:"url"`
}

// NewArticle 先创建文章，再传给NewNews()
func NewArticle(title, desc, picUrl, url string) Article {
	return Article{CDATA(title), CDATA(desc), CDATA(picUrl), CDATA(url)}
}

// MpNews 加密新闻消息，仅企业号支持
type MpNews struct {
	wxResp
	MpNews struct {
		Articles []MpArticle `json:"articles"`
	} `json:"mpnews"`
}

// NewMpNews 加密新闻mpnews消息(仅企业号可用)
func NewMpNews(to string, id int, arts ...MpArticle) MpNews {
	news := MpNews{wxResp: newWxResp(TypeMpNews, to, id)}
	news.MpNews.Articles = arts
	return news
}

// MpArticle 加密文章
type MpArticle struct {
	Title        string `json:"title"`
	ThumbMediaId string `json:"thumb_media_id"`
	Author       string `json:"author"`
	Url          string `json:"content_source_url"`
	Content      string `json:"content"`
	Digest       string `json:"digest"`
	ShowCoverPic int    `json:"show_cover_pic"`
}

// NewMpArticle 先创建加密文章，再传给NewMpNews()
func NewMpArticle(title, mediaId, author, url, content, digest string, showCoverPic int) MpArticle {
	return MpArticle{title, mediaId, author, url, content, digest, showCoverPic}
}
