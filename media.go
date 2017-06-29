package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/esap/wechat/util"
)

// Media 上传回复体
type Media struct {
	WxErr
	Type         string `json:"type"`
	MediaID      string `json:"media_id"`
	ThumbMediaId string `json:"thumb_media_id"`
	CreatedAt    string `json:"created_at"`
}

// MediaUpload 临时素材上传，mediaType可选项：
//	TypeImage  = "image"
//	TypeVoice  = "voice"
//	TypeVideo  = "video"
//	TypeFile   = "file" // 仅企业号可用
//	TypeThumb  = "thumb"
func MediaUpload(mediaType string, filename string) (media Media, err error) {
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

//GetMedia 下载媒体
func GetMedia(filename, mediaId string) error {
	url := fmt.Sprintf(getMedia, GetAccessToken(), mediaId)
	return util.GetFile(filename, url)
}
