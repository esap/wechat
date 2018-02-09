package wechat

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/esap/wechat/util"
)

const (
	// 临时素材上传
	WXAPI_UPLOAD = "media/upload?access_token=%s&type=%s"
	// 临时素材下载
	WXAPI_GETMEDIA = "media/get?access_token=%s&media_id=%s"
	// 高清语言素材下载
	WXAPI_GetJssdkMedia = "media/get/jssdk?access_token=%s&media_id=%s"
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

// GetMedia 下载临时素材
func (s *Server) GetMedia(filename, mediaId string) error {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GETMEDIA, s.GetAccessToken(), mediaId)
	return util.GetFile(filename, url)
}

// GetMedia 下载临时素材
func GetMedia(filename, mediaId string) error {
	return std.GetMedia(filename, mediaId)
}

// GetMediaBytes 下载临时素材,返回body字节
func (s *Server) GetMediaBytes(mediaId string) ([]byte, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GETMEDIA, s.GetAccessToken(), mediaId)
	return util.GetBody(url)
}

// GetMediaBytes 下载媒体,返回io.Reader
func (s *Server) GetBody(mediaId string) (io.ReadCloser, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GETMEDIA, s.GetAccessToken(), mediaId)
	return util.GetRawBody(url)
}

// GetMediaBytes 下载媒体,返回body字节
func GetMediaBytes(mediaId string) ([]byte, error) {
	return std.GetMediaBytes(mediaId)
}

// GetJsMedia 下载高清语言素材(通过JSSDK上传)
func (s *Server) GetJsMedia(filename, mediaId string) error {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GetJssdkMedia, s.GetAccessToken(), mediaId)
	return util.GetFile(filename, url)
}

// GetJsMedia 下载高清语言素材(通过JSSDK上传)
func GetJsMedia(filename, mediaId string) error {
	return std.GetJsMedia(filename, mediaId)
}

// GetJsMediaBytes 下载高清语言素材,返回body字节
func (s *Server) GetJsMediaBytes(mediaId string) ([]byte, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GetJssdkMedia, s.GetAccessToken(), mediaId)
	return util.GetBody(url)
}

// GetJsMediaBody 下载高清语言素材,返回io.Reader
func (s *Server) GetJsMediaBody(mediaId string) (io.ReadCloser, error) {
	url := fmt.Sprintf(s.RootUrl+WXAPI_GetJssdkMedia, s.GetAccessToken(), mediaId)
	return util.GetRawBody(url)
}

// GetJsMediaBytes 下载高清语言素材,返回body字节
func GetJsMediaBytes(mediaId string) ([]byte, error) {
	return std.GetMediaBytes(mediaId)
}
