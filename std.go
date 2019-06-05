package wechat

// import (
// 	"net/http"
// )

// // std 默认单例，未来会完全移除
// var std *Server

// // Set 设置token,appId,secret
// func Set(tk, id, sec string, key ...string) (err error) {
// 	std = NewServer(nil)
// 	return std.Set(tk, id, sec, key...)
// }

// // SetEnt 初始化单例企业微信应用，只有部分接口功能，完整实例需要使用New()
// func SetEnt(token, appId, secret, aeskey string, agentId ...int) (err error) {
// 	return std.SetEnt(token, appId, secret, aeskey, agentId...)
// }

// // VerifyURL 验证URL,验证成功则返回标准请求载体（Msg已解密）
// func VerifyURL(w http.ResponseWriter, r *http.Request) (ctx *Context) {
// 	return std.VerifyURL(w, r)
// }

// // GetAccessToken 读取默认实例AccessToken
// func GetAccessToken() string {
// 	return std.GetAccessToken()
// }

// // GetUserAccessToken 获取默认实例通讯录AccessToken
// func GetUserAccessToken() string {
// 	return std.GetUserAccessToken()
// }

// // GetUserOauth 通过code鉴权
// func GetUserOauth(code string) (userOauth UserOauth, err error) {
// 	return std.GetUserOauth(code)
// }

// // MediaUpload 临时素材上传
// func MediaUpload(mediaType string, filename string) (media Media, err error) {
// 	return std.MediaUpload(mediaType, filename)
// }

// // GetMedia 下载临时素材
// func GetMedia(filename, mediaId string) error {
// 	return std.GetMedia(filename, mediaId)
// }

// // GetMediaBytes 下载媒体,返回body字节
// func GetMediaBytes(mediaId string) ([]byte, error) {
// 	return std.GetMediaBytes(mediaId)
// }

// // GetJsMedia 下载高清语言素材(通过JSSDK上传)
// func GetJsMedia(filename, mediaId string) error {
// 	return std.GetJsMedia(filename, mediaId)
// }

// // GetJsMediaBytes 下载高清语言素材,返回body字节
// func GetJsMediaBytes(mediaId string) ([]byte, error) {
// 	return std.GetMediaBytes(mediaId)
// }

// // SendMsg 发送消息
// func SendMsg(v interface{}) *WxErr {
// 	return std.SendMsg(v)
// }

// // SendText 发送客服text消息
// func SendText(to string, id int, msg string, safe ...int) *WxErr {
// 	return std.SendText(to, id, msg, safe...)
// }

// // SendImage 发送客服Image消息
// func SendImage(to string, id int, mediaId string) *WxErr {
// 	return std.SendImage(to, id, mediaId)
// }

// // SendVoice 发送客服Voice消息
// func SendVoice(to string, id int, mediaId string) *WxErr {
// 	return std.SendVoice(to, id, mediaId)
// }

// // SendFile 发送客服File消息
// func SendFile(to string, id int, mediaId string) *WxErr {
// 	return std.SendFile(to, id, mediaId)
// }

// // SendVideo 发送客服Video消息
// func SendVideo(to string, id int, mediaId, title, desc string) *WxErr {
// 	return std.SendVideo(to, id, mediaId, title, desc)
// }

// // SendTextcard 发送客服extcard消息
// func SendTextcard(to string, id int, title, desc, url string) *WxErr {
// 	return std.SendTextcard(to, id, title, desc, url)
// }

// // SendMusic 发送客服Music消息
// func SendMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) *WxErr {
// 	return std.SendMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl)
// }

// // SendNews 发送客服news消息
// func SendNews(to string, id int, arts ...Article) *WxErr {
// 	return std.SendNews(to, id, arts...)
// }

// // SendMpNews 发送加密新闻mpnews消息(仅企业号可用)
// func SendMpNews(to string, id int, arts ...MpArticle) *WxErr {
// 	return std.SendMpNews(to, id, arts...)
// }

// // SendMpNewsId 发送加密新闻mpnews消息(直接使用mediaId)
// func SendMpNewsId(to string, id int, mediaId string) *WxErr {
// 	return std.SendMpNewsId(to, id, mediaId)
// }

// func newWxResp(msgType, toUser string, agentId int) wxResp {
// 	return std.newWxResp(msgType, toUser, agentId)
// }

// // NewText Text 文本消息
// func NewText(to string, id int, msg ...string) Text {
// 	return std.NewText(to, id, msg...)
// }

// // NewImage Image 消息
// func NewImage(to string, id int, mediaId string) Image {
// 	return std.NewImage(to, id, mediaId)
// }

// // NewVoice Voice消息
// func NewVoice(to string, id int, mediaId string) Voice {
// 	return std.NewVoice(to, id, mediaId)
// }

// // NewFile File消息
// func NewFile(to string, id int, mediaId string) File {
// 	return std.NewFile(to, id, mediaId)
// }

// // NewVideo Video消息
// func NewVideo(to string, id int, mediaId, title, desc string) Video {
// 	return std.NewVideo(to, id, mediaId, title, desc)
// }

// // NewTextcard Textcard消息
// func NewTextcard(to string, id int, title, description, url string) Textcard {
// 	return std.NewTextcard(to, id, title, description, url)
// }

// // NewMusic Music消息
// func NewMusic(to string, id int, mediaId, title, desc, musicUrl, qhMusicUrl string) Music {
// 	return std.NewMusic(to, id, mediaId, title, desc, musicUrl, qhMusicUrl)
// }

// // NewNews news消息
// func NewNews(to string, id int, arts ...Article) (news News) {
// 	return std.NewNews(to, id, arts...)
// }

// // NewMpNews 加密新闻mpnews消息(仅企业微信可用)
// func NewMpNews(to string, id int, arts ...MpArticle) (news MpNews) {
// 	return std.NewMpNews(to, id, arts...)
// }

// // NewMpNewsId 加密新闻mpnews消息(仅企业微信可用)
// func NewMpNewsId(to string, id int, mediaId string) (news MpNewsId) {
// 	return std.NewMpNewsId(to, id, mediaId)
// }

// // NewWxCard 卡券消息，服务号可用
// func NewWxCard(to string, id int, cardId string) WxCard {
// 	return std.NewWxCard(to, id, cardId)
// }

// // NewMarkDown markdown消息，企业微信可用
// func NewMarkDown(to string, id int, content string) MarkDown {
// 	return std.NewMarkDown(to, id, content)
// }
