package wechat

import (
	"encoding/xml"
)

// WxMsg 混合用户消息，业务判断的主体
type WxMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string // text
	MsgId        int64
	AgentID      int     // corp
	PicUrl       string  // image
	MediaId      string  // image/voice/video/shortvideo
	Format       string  // voice
	Recognition  string  // voice
	ThumbMediaId string  // video
	LocationX    float32 `xml:"Latitude"`  // location
	LocationY    float32 `xml:"Longitude"` // location
	Precision    float32 // LOCATION
	Scale        int     // location
	Label        string  // location
	Title        string  // link
	Description  string  // link
	Url          string  // link
	Event        string  // event
	EventKey     string  // event
	Ticket       string
	ScanCodeInfo struct {
		ScanType   string
		ScanResult string
	}
}

// WxMsgEnc 密文用户消息
type WxMsgEnc struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	AgentID    int
	Encrypt    string
}
