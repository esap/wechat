package wechat

import (
	"encoding/xml"
)

type (
	// WxMsg 混合用户消息，业务判断的主体
	WxMsg struct {
		XMLName      xml.Name `xml:"xml"`
		ToUserName   string
		FromUserName string
		CreateTime   int64
		MsgId        int64
		MsgType      string
		Content      string  // text
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

		AgentType string
		ItemCount int
		PackageId string

		Item []struct {
			FromUserName string
			CreateTime   int64
			MsgType      string
			Event        string // event
			Name         string
			Owner        string
			AddUserList  string
			DelUserList  string
			ChatId       string
			MsgId        int64

			ChatInfo struct {
				ChatId   string
				Name     string
				Owner    string
				UserList string
			}

			Content  string // text
			Receiver struct {
				Type string
				Id   string
			}

			FileName string // file
			PicUrl   string // image/link
			MediaId  string // image/voice/video/shortvideo

			Location_X float32 // location
			Location_Y float32 // location
			Scale      int     // location
			Label      string  // location

			Title       string // link
			Description string // link
			Url         string // link
		}
	}

	// WxMsgEnc 加密的用户消息
	WxMsgEnc struct {
		XMLName    xml.Name `xml:"xml"`
		ToUserName string
		AgentID    int
		Encrypt    string
		AgentType  string
	}
)
