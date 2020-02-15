// 目前官方未提供golang版，本SDK实现参考了php版官方库
// @woylin, since 2016-1-6

package wechat

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/esap/wechat/util"
)

// WXAPI 订阅号，服务号，小程序接口，相关接口常量统一以此开头
const (
	WXAPI      = "https://api.weixin.qq.com/cgi-bin/"
	WXAPIToken = WXAPI + "token?grant_type=client_credential&appid=%s&secret=%s"
	WXAPIMsg   = WXAPI + "message/custom/send?access_token="
	WXAPIJsapi = WXAPI + "get_jsapi_ticket?access_token="
)

// CorpAPI 企业微信接口，相关接口常量统一以此开头
const (
	CorpAPI      = "https://qyapi.weixin.qq.com/cgi-bin/"
	CorpAPIToken = CorpAPI + "gettoken?corpid=%s&corpsecret=%s"
	CorpAPIMsg   = CorpAPI + "message/send?access_token="
	CorpAPIJsapi = CorpAPI + "get_jsapi_ticket?access_token="
)

var (
	// Debug is a flag to Println()
	Debug bool = false

	// UserServerMap 通讯录实例集，用于企业微信
	UserServerMap = make(map[string]*Server)
)

// WxConfig 配置，用于New()
type WxConfig struct {
	AppId                string
	Token                string
	Secret               string
	EncodingAESKey       string
	AgentId              int
	MchId                string
	AppName              string
	AppType              int                                  // 0-公众号,小程序; 1-企业微信
	ExternalTokenHandler func(string, ...string) *AccessToken // 外部token获取函数
}

// Server 微信服务容器
type Server struct {
	AppId   string
	MchId   string // 商户id，用于微信支付
	AgentId int
	Secret  string

	Token          string
	EncodingAESKey string

	AppName  string // 唯一标识，主要用于企业微信多应用区分
	AppType  int    // 0-公众号,小程序; 1-企业微信
	AesKey   []byte // 解密的AesKey
	SafeMode bool
	EntMode  bool

	RootUrl  string
	MsgUrl   string
	TokenUrl string
	JsApi    string

	Safe        int
	accessToken *AccessToken
	ticket      *Ticket
	UserList    userList
	DeptList    DeptList
	TagList     TagList
	MsgQueue    chan interface{}
	sync.Mutex  // accessToken读取锁

	ExternalTokenHandler func(appId string, appName ...string) *AccessToken // 通过外部方法统一获取access token ,避免集群情况下token失效
}

func Set(wc *WxConfig) *Server {
        return &Server{
                AppId:                wc.AppId,
                Secret:               wc.Secret,
                AgentId:              wc.AgentId,
                MchId:                wc.MchId,
                AppName:              wc.AppName,
                AppType:              wc.AppType,
                Token:                wc.Token,
                EncodingAESKey:       wc.EncodingAESKey,
                ExternalTokenHandler: wc.ExternalTokenHandler,
        }
}
// New 微信服务容器
func New(wc *WxConfig) *Server {
	s := Set(wc)

	switch wc.AppType {
	case 1:
		s.RootUrl = CorpAPI
		s.MsgUrl = CorpAPIMsg
		s.TokenUrl = CorpAPIToken
		s.JsApi = CorpAPIJsapi
		s.EntMode = true
	default:
		s.RootUrl = WXAPI
		s.MsgUrl = WXAPIMsg
		s.TokenUrl = WXAPIToken
		s.JsApi = WXAPIJsapi
	}

	err := s.getAccessToken()
	if err != nil {
		log.Println("getAccessToken err:", err)
	}

	// 存在EncodingAESKey则开启加密安全模式
	if len(s.EncodingAESKey) > 0 && s.EncodingAESKey != "" {
		s.SafeMode = true
		if s.AesKey, err = base64.StdEncoding.DecodeString(s.EncodingAESKey + "="); err != nil {
			log.Println("AesKey解析错误:", err)
		}
		Println("启用加密模式")
	}

	if s.AgentId == 9999999 {
		UserServerMap[s.AppId] = s // 这里约定传入企业微信通讯录secret时，agentId=9999999
	}

	if s.AppType == 1 {
		s.FetchUserList()
	}

	s.MsgQueue = make(chan interface{}, 1000)
	go func() {
		for {
			msg := <-s.MsgQueue
			e := s.SendMsg(msg)
			if e.ErrCode != 0 {
				log.Println("MsgSend err:", e.ErrMsg)
			}
		}
	}()

	return s
}

// VerifyURL 验证URL,验证成功则返回标准请求载体（Msg已解密）
func (s *Server) VerifyURL(w http.ResponseWriter, r *http.Request) (ctx *Context) {
	Println(r.Method, "|", r.URL.String())
	ctx = &Context{
		Server:    s,
		Writer:    w,
		Request:   r,
		Timestamp: r.FormValue("timestamp"),
		Nonce:     r.FormValue("nonce"),
		Msg:       new(WxMsg),
	}

	// 明文模式可直接解析body->消息
	if !s.SafeMode && r.Method == "POST" {
		if err := xml.NewDecoder(r.Body).Decode(ctx.Msg); err != nil {
			Println("Decode WxMsg err:", err)
		}
	}

	// 密文模式，消息在body.Encrypt
	echostr := r.FormValue("echostr")
	if s.SafeMode && r.Method == "POST" {
		msgEnc := new(WxMsgEnc)
		if err := xml.NewDecoder(r.Body).Decode(msgEnc); err != nil {
			Println("Decode MsgEnc err:", err)
		}
		echostr = msgEnc.Encrypt
	}

	// 验证signature
	signature := r.FormValue("signature")
	if signature == "" {
		signature = r.FormValue("msg_signature")
	}
	if s.EntMode && signature != util.SortSha1(s.Token, ctx.Timestamp, ctx.Nonce, echostr) {
		log.Println("Signature验证错误!(企业微信)", s.Token, ctx.Timestamp, ctx.Nonce, echostr)
		return
	} else if !s.EntMode && signature != util.SortSha1(s.Token, ctx.Timestamp, ctx.Nonce) {
		log.Println("Signature验证错误!(公众号)", util.SortSha1(s.Token, ctx.Timestamp, ctx.Nonce))
		return
	}

	// 密文模式，解密echostr中的消息
	if s.EntMode || (s.SafeMode && r.Method == "POST") {
		var err error
		echostr, err = s.DecryptMsg(echostr)
		if err != nil {
			log.Println("DecryptMsg error:", err)
			return
		}
	}

	if r.Method == "GET" {
		Println("api echostr:", echostr)
		w.Write([]byte(echostr))
		return
	}

	Println("Wechat ==>", echostr)
	if s.SafeMode {
		if err := xml.Unmarshal([]byte(echostr), ctx.Msg); err != nil {
			log.Println("Msg parse err:", err)
		}
	}

	return
}

// DecryptMsg 解密微信消息,密文string->base64Dec->aesDec->去除头部随机字串
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppId组成
func (s *Server) DecryptMsg(msg string) (string, error) {
	aesMsg, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	buf, err := util.AesDecrypt(aesMsg, s.AesKey)
	if err != nil {
		return "", err
	}

	var msgLen int32
	binary.Read(bytes.NewBuffer(buf[16:20]), binary.BigEndian, &msgLen)
	if msgLen < 0 || msgLen > 1000000 {
		return "", errors.New("AesKey is invalid")
	}
	if string(buf[20+msgLen:]) != s.AppId {
		return "", errors.New("AppId is invalid")
	}
	return string(buf[20 : 20+msgLen]), nil
}

// wxRespEnc 加密回复体
type wxRespEnc struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA
	MsgSignature CDATA
	TimeStamp    string
	Nonce        CDATA
}

// EncryptMsg 加密普通回复(AES-CBC),打包成xml格式
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppId组成
func (s *Server) EncryptMsg(msg []byte, timeStamp, nonce string) (re *wxRespEnc, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(msg)))
	if err != nil {
		return
	}
	l := buf.Bytes()

	rd := []byte(util.GetRandomString(16))

	plain := bytes.Join([][]byte{rd, l, msg, []byte(s.AppId)}, nil)
	ae, _ := util.AesEncrypt(plain, s.AesKey)
	encMsg := base64.StdEncoding.EncodeToString(ae)
	re = &wxRespEnc{
		Encrypt:      CDATA(encMsg),
		MsgSignature: CDATA(util.SortSha1(s.Token, timeStamp, nonce, encMsg)),
		TimeStamp:    timeStamp,
		Nonce:        CDATA(nonce),
	}
	return
}

// SetLog 设置log
func SetLog(l io.Writer) {
	log.SetOutput(l)
}

// SafeOpen 设置密保模式
func (s *Server) SafeOpen() {
	s.Safe = 1
}

// SafeClose 关闭密保模式
func (s *Server) SafeClose() {
	s.Safe = 0
}

// Println Debug输出
func Println(v ...interface{}) {
	if Debug {
		log.Println(v...)
	}
}

// Printf Debug输出
func Printf(s string, v ...interface{}) {
	if Debug {
		log.Printf(s, v...)
	}
}
