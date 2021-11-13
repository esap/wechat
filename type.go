package wechat

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// Type io类型汇总
const (
	TypeText     = "text"
	TypeImage    = "image"
	TypeVoice    = "voice"
	TypeMusic    = "music"
	TypeVideo    = "video"
	TypeTextcard = "textcard" // 仅企业微信可用
	TypeWxCard   = "wxcard"   // 仅服务号可用
	TypeMarkDown = "markdown" // 仅企业微信可用
	TypeTaskCard = "taskcard" // 仅企业微信可用
	TypeFile     = "file"     // 仅企业微信可用
	TypeNews     = "news"
	TypeMpNews   = "mpnews" // 仅企业微信可用
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

// wxResp 响应消息共用字段
// 响应消息被动回复为XML结构，文本类型采用CDATA编码规范
// 响应消息主动发送为json结构，即客服消息
type wxResp struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   CDATA    `json:"touser"`
	ToParty      CDATA    `xml:"-" json:"toparty"` // 企业号专用
	ToTag        CDATA    `xml:"-" json:"totag"`   // 企业号专用
	FromUserName CDATA    `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      CDATA    `json:"msgtype"`
	AgentId      int      `xml:"-" json:"agentid"`
	Safe         int      `xml:"-" json:"safe"`
}

// to字段格式："userid1|userid2 deptid1|deptid2 tagid1|tagid2"
func (s *Server) newWxResp(msgType, to string) (r wxResp) {
	toArr := strings.Split(to, " ")
	r = wxResp{
		ToUserName: CDATA(toArr[0]),
		MsgType:    CDATA(msgType),
		AgentId:    s.AgentId,
		Safe:       s.Safe}
	if len(toArr) > 1 {
		r.ToParty = CDATA(toArr[1])
	}
	if len(toArr) > 2 {
		r.ToTag = CDATA(toArr[2])
	}
	return
}

// Text 文本消息
type (
	Text struct {
		wxResp
		content `xml:"Content" json:"text"`
	}

	content struct {
		Content CDATA `json:"content"`
	}
)

// NewText Text 文本消息
func (s *Server) NewText(to string, msg ...string) Text {
	return Text{
		s.newWxResp(TypeText, to),
		content{CDATA(strings.Join(msg, ""))},
	}
}

// Image 图片消息
type (
	Image struct {
		wxResp
		Image media `json:"image"`
	}

	media struct {
		MediaId CDATA `json:"media_id"`
	}
)

// NewImage Image 消息
func (s *Server) NewImage(to, mediaId string) Image {
	return Image{
		s.newWxResp(TypeImage, to),
		media{CDATA(mediaId)},
	}
}

// Voice 语音消息
type Voice struct {
	wxResp
	Voice media `json:"voice"`
}

// NewVoice Voice消息
func (s *Server) NewVoice(to, mediaId string) Voice {
	return Voice{
		s.newWxResp(TypeVoice, to),
		media{CDATA(mediaId)},
	}
}

// File 文件消息，仅企业号支持
type File struct {
	wxResp
	File media `json:"file"`
}

// NewFile File消息
func (s *Server) NewFile(to, mediaId string) File {
	return File{
		s.newWxResp(TypeFile, to),
		media{CDATA(mediaId)},
	}
}

// Video 视频消息
type (
	Video struct {
		wxResp
		Video video `json:"video"`
	}

	video struct {
		MediaId     CDATA `json:"media_id"`
		Title       CDATA `json:"title"`
		Description CDATA `json:"description"`
	}
)

// NewVideo Video消息
func (s *Server) NewVideo(to, mediaId, title, desc string) Video {
	return Video{
		s.newWxResp(TypeVideo, to),
		video{CDATA(mediaId), CDATA(title), CDATA(desc)},
	}
}

// Textcard 卡片消息，仅企业微信客户端有效
type (
	Textcard struct {
		wxResp
		Textcard textcard `json:"textcard"`
	}

	textcard struct {
		Title       CDATA `json:"title"`
		Description CDATA `json:"description"`
		Url         CDATA `json:"url"`
	}
)

// NewTextcard Textcard消息
func (s *Server) NewTextcard(to, title, description, url string) Textcard {
	return Textcard{
		s.newWxResp(TypeTextcard, to),
		textcard{CDATA(title), CDATA(description), CDATA(url)},
	}
}

// Music 音乐消息，企业微信不支持
type (
	Music struct {
		wxResp
		Music music `json:"music"`
	}

	music struct {
		Title        CDATA `json:"title"`
		Description  CDATA `json:"description"`
		MusicUrl     CDATA `json:"musicurl"`
		HQMusicUrl   CDATA `json:"hqmusicurl"`
		ThumbMediaId CDATA `json:"thumb_media_id"`
	}
)

// NewMusic Music消息
func (s *Server) NewMusic(to, mediaId, title, desc, musicUrl, qhMusicUrl string) Music {
	return Music{
		s.newWxResp(TypeMusic, to),
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
func (s *Server) NewNews(to string, arts ...Article) (news News) {
	news.wxResp = s.newWxResp(TypeNews, to)
	news.ArticleCount = len(arts)
	news.Articles.Item = arts
	return
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

type (
	// MpNews 加密新闻消息，仅企业微信支持
	MpNews struct {
		wxResp
		MpNews struct {
			Articles []MpArticle `json:"articles"`
		} `json:"mpnews"`
	}

	// MpNewsId 加密新闻消息(通过mediaId直接发)
	MpNewsId struct {
		wxResp
		MpNews struct {
			MediaId CDATA `json:"media_id"`
		} `json:"mpnews"`
	}
)

// NewMpNews 加密新闻mpnews消息(仅企业微信可用)
func (s *Server) NewMpNews(to string, arts ...MpArticle) (news MpNews) {
	news.wxResp = s.newWxResp(TypeMpNews, to)
	news.MpNews.Articles = arts
	return
}

// NewMpNewsId 加密新闻mpnews消息(仅企业微信可用)
func (s *Server) NewMpNewsId(to string, mediaId string) (news MpNewsId) {
	news.wxResp = s.newWxResp(TypeMpNews, to)
	news.MpNews.MediaId = CDATA(mediaId)
	return
}

// MpArticle 加密文章
type MpArticle struct {
	Title        string `json:"title"`
	ThumbMediaId string `json:"thumb_media_id"`
	Author       string `json:"author"`
	Url          string `json:"content_source_url"`
	Content      string `json:"content"`
	Digest       string `json:"digest"`
}

// NewMpArticle 先创建加密文章，再传给NewMpNews()
func NewMpArticle(title, mediaId, author, url, content, digest string) MpArticle {
	return MpArticle{title, mediaId, author, url, content, digest}
}

// WxCard 卡券
type WxCard struct {
	wxResp
	WxCard struct {
		CardId string `json:"card_id"`
	} `json:"wxcard"`
}

// NewWxCard 卡券消息，服务号可用
func (s *Server) NewWxCard(to, cardId string) (c WxCard) {
	c.wxResp = s.newWxResp(TypeWxCard, to)
	c.WxCard.CardId = cardId
	return
}

// MarkDown markdown消息，仅企业微信支持，上限2048字节，utf-8编码
type MarkDown struct {
	wxResp
	MarkDown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

// NewMarkDown markdown消息，企业微信可用
func (s *Server) NewMarkDown(to, content string) (md MarkDown) {
	md.wxResp = s.newWxResp(TypeMarkDown, to)
	md.MarkDown.Content = content
	return
}

// TaskCard 任务卡片消息，仅企业微信支持，支持一到两个按钮设置
type TaskCard struct {
	wxResp
	TaskCard struct {
		Title       string                   `json:"title"`
		Description string                   `json:"description"`
		Url         string                   `json:"url"`
		TaskId      string                   `json:"task_id"`
		Btn         []map[string]interface{} `json:"btn"`
	} `json:"taskcard"`
}

// NewTaskCard 任务卡片消息，企业微信可用
func (s *Server) NewTaskCard(to, Title, Desc, Url, TaskId, Btn string) (tc TaskCard) {
	tc.wxResp = s.newWxResp(TypeTaskCard, to)
	tc.TaskCard.Title = Title
	tc.TaskCard.Description = Desc
	tc.TaskCard.Url = Url
	tc.TaskCard.TaskId = TaskId
	mp := make([]map[string]interface{}, 0)
	if Btn != "" {
		if err := json.Unmarshal([]byte(Btn), &mp); err != nil {
			fmt.Println("create taskcard btn err:", err)
		} else {
			tc.TaskCard.Btn = mp
		}
	}
	return
}
