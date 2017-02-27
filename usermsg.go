package wechat

import (
	"encoding/xml"
	"log"
	"net/http"

	"github.com/esap/wechat/util"
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

// parseWxReq 解析http请求，返回请求体
func parseWxMsg(r *http.Request) (msg *WxMsg) {
	msg = new(WxMsg)
	if err := util.HttpParseXml(r, msg); err != nil {
		log.Println("parseWxMsg err:", err)
	}
	Println("ReqXml=\n", msg)
	return
}

// wxMsgEnc 密文用户消息
type wxMsgEnc struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	AgentID    string
	Encrypt    string
}

// parseEncReq 解析加密请求XML
func parseEncMsg(r *http.Request) (msg *wxMsgEnc) {
	msg = new(wxMsgEnc)
	if err := util.HttpParseXml(r, msg); err != nil {
		log.Println("parseEncMsg err:", err)
	}
	return
}
