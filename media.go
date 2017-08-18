package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/esap/wechat/util"
)

const (
	WXAPI_UPLOAD   = "media/upload?access_token=%s&type=%s"
	WXAPI_GETMEDIA = "media/get?access_token=%s&media_id=%s"
)

// Media 上传回复体
type Media struct {
	WxErr
	Type         string      `json:"type"`
	MediaID      string      `json:"media_id"`
	ThumbMediaId string      `json:"thumb_media_id"`
	CreatedAt    interface{} `json:"created_at"` // 企业号是string,服务号是int,采用interface{}统一接收
}

// MediaUpload 临时素材上传，mediaType选项如下：
//	TypeImage  = "image"
//	TypeVoice  = "voice"
//	TypeVideo  = "video"
//	TypeFile   = "file" // 仅企业号可用
//	TypeThumb  = "thumb"
func (s *Server) MediaUpload(mediaType string, filename string) (media Media, err error) {
	uri := fmt.Sprintf(s.RootUrl+WXAPI_UPLOAD, s.GetAccessToken(), mediaType)
	var b []byte
	b, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &media); err != nil {
		return
	}
	err = media.Error()
	return
}

// MediaUpload 临时素材上传
func MediaUpload(mediaType string, filename string) (media Media, err error) {
	return std.MediaUpload(mediaType, filename)
}

// GetMedia 下载媒体
func (s *Server) GetMedia(filename, mediaId string) error {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GETMEDIA, s.GetAccessToken(), mediaId)
	return util.GetFile(filename, url)
}

// GetMedia 下载媒体
func GetMedia(filename, mediaId string) error {
	return std.GetMedia(filename, mediaId)
}

// GetMediaBytes 下载媒体
func (s *Server) GetMediaBytes(mediaId string) ([]byte, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GETMEDIA, s.GetAccessToken(), mediaId)
	return util.GetBody(url)
}

// GetMediaBytes 下载媒体
func GetMediaBytes(mediaId string) ([]byte, error) {
	return std.GetMediaBytes(mediaId)
}
