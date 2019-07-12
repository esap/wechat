package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/esap/wechat/util"
)

const (
	// WXAPIMediaUpload 临时素材上传
	WXAPIMediaUpload = "media/upload?access_token=%s&type=%s"
	// WXAPIMediaGet 临时素材下载
	WXAPIMediaGet = "media/get?access_token=%s&media_id=%s"
	// WXAPIMediaGetJssdk 高清语言素材下载
	WXAPIMediaGetJssdk = "media/get/jssdk?access_token=%s&media_id=%s"
)

// Media 上传回复体
type Media struct {
	WxErr
	Type         string      `json:"type"`
	MediaID      string      `json:"media_id"`
	ThumbMediaId string      `json:"thumb_media_id"`
	CreatedAt    interface{} `json:"created_at"` // 企业微信是string,服务号是int,采用interface{}统一接收
}

// MediaUpload 临时素材上传，mediaType选项如下：
//	TypeImage  = "image"
//	TypeVoice  = "voice"
//	TypeVideo  = "video"
//	TypeFile   = "file" // 仅企业微信可用
func (s *Server) MediaUpload(mediaType string, filename string) (media Media, err error) {
	uri := fmt.Sprintf(s.RootUrl+WXAPIMediaUpload, s.GetAccessToken(), mediaType)
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

// GetMedia 下载临时素材
func (s *Server) GetMedia(filename, mediaId string) error {
	url := fmt.Sprintf(s.RootUrl+WXAPIMediaGet, s.GetAccessToken(), mediaId)
	return util.GetFile(filename, url)
}

// GetMediaBytes 下载临时素材,返回body字节
func (s *Server) GetMediaBytes(mediaId string) ([]byte, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPIMediaGet, s.GetAccessToken(), mediaId)
	return util.GetBody(url)
}

// GetJsMedia 下载高清语言素材(通过JSSDK上传)
func (s *Server) GetJsMedia(filename, mediaId string) error {
	url := fmt.Sprintf(s.RootUrl+WXAPIMediaGetJssdk, s.GetAccessToken(), mediaId)
	return util.GetFile(filename, url)
}

// GetJsMediaBytes 下载高清语言素材,返回body字节
func (s *Server) GetJsMediaBytes(mediaId string) ([]byte, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPIMediaGetJssdk, s.GetAccessToken(), mediaId)
	return util.GetBody(url)
}
