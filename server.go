// 目前官方未提供golang版，本SDK实现参考了php版官方库
// @woylin, since 2016-1-6

package wechat

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 订阅号，服务号接口
const (
	WXAPI       = "https://api.weixin.qq.com/cgi-bin/"
	WXAPI_TOKEN = WXAPI + "token?grant_type=client_credential&appid=%s&secret=%s"
	WXAPI_MSG   = WXAPI + "message/custom/send?access_token="
)

var (
	// Debug is a flag to Println()
	Debug bool = false
	std        = NewServer()
)

// Server 微信服务容器
type Server struct {
	AppId          string
	AgentId        int
	Secret         string
	Token          string
	EncodingAESKey string
	AesKey         []byte // 解密的AesKey
	SafeMode       bool
	EntMode        bool
	RootUrl        string
	MsgUrl         string
	TokenUrl       string
	Safe           int
	accessToken    *AccessToken
	UserList       userList
	DeptList       DeptList
	TagList        TagList
	MsgQueue       chan interface{}
}

// New 微信服务容器，根据agentId判断是企业号或服务号
func New(token, appid, secret, key string, agentId ...int) (s *Server) {
	s = NewServer()
	if len(agentId) > 0 {
		s.SetEnt(token, appid, secret, key, agentId[0])
	} else {
		s.Set(token, appid, secret, key)
	}
	return s
}

// NewServer 空容器
func NewServer() *Server {
	s := &Server{
		RootUrl:  WXAPI,
		MsgUrl:   WXAPI_MSG,
		TokenUrl: WXAPI_TOKEN,
	}
	s.init()
	return s
}

//SetLog 设置log
func SetLog(l io.Writer) {
	log.SetOutput(l)
}

//SafeOpen 设置密保模式
func (s *Server) SafeOpen() {
	s.Safe = 1
}

//SafeOpen 设置密保模式
func (s *Server) SafeClose() {
	s.Safe = 0
}

// Set 设置token,appId,secret
func (s *Server) Set(tk, id, sec string, key ...string) (err error) {
	s.Token, s.AppId, s.Secret = tk, id, sec
	// 存在EncodingAESKey则开启加密安全模式
	if len(key) > 0 {
		s.SafeMode = true
		if s.AesKey, err = base64.StdEncoding.DecodeString(key[0] + "="); err != nil {
			return err
		}
		Println("启用加密模式")
	}
	return
}

// Set 设置token,appId,secret
func Set(tk, id, sec string, key ...string) (err error) {
	return std.Set(tk, id, sec, key...)
}

// VerifyURL 验证URL,验证成功则返回标准请求载体（Msg已解密）
func (s *Server) VerifyURL(w http.ResponseWriter, r *http.Request) (ctx *Context) {
	Println(r.Method, "|", r.URL.String())
	ctx = &Context{
		Server:    s,
		Writer:    w,
		Request:   r,
		repCount:  0,
		Timestamp: r.FormValue("timestamp"),
		Nonce:     r.FormValue("nonce"),
		Msg:       new(WxMsg),
		MsgEnc:    new(WxMsgEnc),
	}
	if !s.SafeMode && r.Method == "POST" {
		if err := xml.NewDecoder(r.Body).Decode(ctx.Msg); err != nil {
			Println("Decode WxMsg err:", err)
		}
	}

	echostr := r.FormValue("echostr")
	if s.SafeMode && r.Method == "POST" {
		if err := xml.NewDecoder(r.Body).Decode(ctx.MsgEnc); err != nil {
			Println("Decode MsgEnc err:", err)
		}
		echostr = ctx.MsgEnc.Encrypt // POST请求需解析消息体中的Encrypt
	}

	// 验证signature
	signature := r.FormValue("signature") + r.FormValue("msg_signature")
	if s.EntMode && signature != sortSha1(s.Token, ctx.Timestamp, ctx.Nonce, echostr) {
		log.Println("Signature验证错误!(企业号)", s.Token, ctx.Timestamp, ctx.Nonce, echostr)
		return
	} else if !s.EntMode && r.FormValue("signature") != sortSha1(s.Token, ctx.Timestamp, ctx.Nonce) {
		log.Println("Signature验证错误!(公众号)", s.Token, ctx.Timestamp, ctx.Nonce)
		return
	}
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
	if r.Method == "POST" && ctx.Msg != nil && ctx.Msg.AgentType == "chat" {
		Println("Chat echostr:", ctx.Msg.PackageId)
		w.Write([]byte(ctx.Msg.PackageId))
	}
	return
}

// VerifyURL 验证URL,验证成功则返回标准请求载体（Msg已解密）
func VerifyURL(w http.ResponseWriter, r *http.Request) (ctx *Context) {
	return std.VerifyURL(w, r)
}

// sortSha1 排序并sha1，主要用于计算signature
func sortSha1(s ...string) string {
	sort.Strings(s)
	h := sha1.New()
	h.Write([]byte(strings.Join(s, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
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

	rd := []byte("welcometoesapsys")

	plain := bytes.Join([][]byte{rd, l, msg, []byte(s.AppId)}, nil)
	ae, _ := util.AesEncrypt(plain, s.AesKey)
	encMsg := base64.StdEncoding.EncodeToString(ae)
	re = &wxRespEnc{
		Encrypt:      CDATA(encMsg),
		MsgSignature: CDATA(sortSha1(s.Token, timeStamp, nonce, encMsg)),
		TimeStamp:    timeStamp,
		Nonce:        CDATA(nonce),
	}
	return
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
