package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/esap/wechat/util"
)

// MediaType 媒体文件类型
// From: https://github.com/silenceper/wechat/material/media.go
type MediaType string

const (
	// MediaTypeImage 媒体文件:图片
	MediaTypeImage MediaType = "image"
	// MediaTypeVoice 媒体文件:声音
	MediaTypeVoice = "voice"
	// MediaTypeVideo 媒体文件:视频
	MediaTypeVideo = "video"
	// MediaTypeFile 媒体文件:文件(企业号可用)
	MediaTypeFile = "file"
	// MediaTypeThumb 媒体文件:缩略图
	MediaTypeThumb = "thumb"
)

// Media 上传回复体
type Media struct {
	WxErr
	Type         MediaType `json:"type"`
	MediaID      string    `json:"media_id"`
	ThumbMediaId string    `json:"thumb_media_id"`
	CreatedAt    int64     `json:"created_at"`
}

// MediaUpload 临时素材上传
func MediaUpload(mediaType MediaType, filename string) (media Media, err error) {
	uri := fmt.Sprintf(uploadUrl, GetAccessToken(), mediaType)
	var b []byte
	b, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &media); err != nil {
		return
	}
	if media.ErrCode != 0 {
		err = fmt.Errorf("MediaUpload error : errcode=%v , errmsg=%v", media.ErrCode, media.ErrMsg)
	}
	return
}
