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
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 订阅号，服务号接口
const (
	WXAPI        = "https://api.weixin.qq.com/cgi-bin/"
	WXAPI_TOKEN  = WXAPI + "token?grant_type=client_credential&appid=%s&secret=%s"
	WXAPI_MSG    = WXAPI + "message/custom/send?access_token="
	WXAPI_UPLOAD = WXAPI + "media/upload?access_token=%s&type=%s"
)

var (
	token     string // 默认token
	appId     string // 企业号填corpId
	secret    string // 管理连接密钥
	aesKey    []byte // 解密的AesKey
	safeMode  bool   = false
	entMode   bool   = false
	msgUrl    string = WXAPI_MSG
	uploadUrl string = WXAPI_UPLOAD
	// Debug is a flag to Println()
	Debug bool = false
)

// Set 设置token,appId,secret
func Set(tk, id, sec string, key ...string) (err error) {
	token = tk
	appId = id
	secret = sec
	// 存在aesKey则开启加密安全模式
	if len(key) > 0 {
		safeMode = true
		aesKey, err = base64.StdEncoding.DecodeString(key[0] + "=")
		if err != nil {
			return err
		}
		log.Println("已启动加密模式")
	}
	FetchAccessToken(WXAPI_TOKEN)
	return nil
}

// VerifyURL 验证URL,验证成功则返回标准请求载体（Msg已解密）
func VerifyURL(w http.ResponseWriter, r *http.Request) (ctx *Context) {
	log.Println(r.Method, "|", r.URL.String())
	ctx = new(Context)
	ctx.Writer = w
	ctx.Request = r
	ctx.repCount = 0
	ctx.Timestamp = r.FormValue("timestamp")
	ctx.Nonce = r.FormValue("nonce")
	signature := r.FormValue("signature")
	if entMode {
		signature = r.FormValue("msg_signature")
	}

	echostr := r.FormValue("echostr")
	if safeMode && r.Method == "POST" {
		echostr = parseEncMsg(r).Encrypt // POST请求需解析消息体中的Encrypt
	}

	// 验证signature
	if entMode && signature != sortSha1(token, ctx.Timestamp, ctx.Nonce, echostr) {
		log.Println("Signature验证错误!(企业号)")
		return
	} else if !entMode && signature != sortSha1(token, ctx.Timestamp, ctx.Nonce) {
		log.Println("Signature验证错误!(公众号)")
		return
	}
	if entMode || (safeMode && r.Method == "POST") {
		var err error
		echostr, err = DecryptMsg(echostr)
		if err != nil {
			log.Println("DecryptMsg error:", err)
			return
		}
	}
	if r.Method == "GET" {
		Println("write echostr:", echostr)
		w.Write([]byte(echostr))
		return
	}
	Println("--ReqEchostr:\n", echostr)
	ctx.Msg = parseWxMsg(r)
	if safeMode {
		body := []byte(echostr)
		if err := xml.Unmarshal(body, ctx.Msg); err != nil {
			log.Println("Msg parse err:", err)
		}
	}
	return
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
func DecryptMsg(s string) (string, error) {
	aesMsg, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	buf, err := util.AesDecrypt(aesMsg, aesKey)
	if err != nil {
		return "", err
	}

	var msgLen int32
	binary.Read(bytes.NewBuffer(buf[16:20]), binary.BigEndian, &msgLen)
	if msgLen < 0 || msgLen > 1000000 {
		return "", errors.New("AesKey is invalid")
	}
	if string(buf[20+msgLen:]) != appId {
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
func EncryptMsg(msg []byte, timeStamp, nonce string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(msg)))
	if err != nil {
		return nil, err
	}
	l := buf.Bytes()

	rd := []byte("welcometoesapsys")

	plain := bytes.Join([][]byte{rd, l, msg, []byte(appId)}, nil)
	ae, _ := util.AesEncrypt(plain, aesKey)
	encMsg := base64.StdEncoding.EncodeToString(ae)
	re := &wxRespEnc{
		Encrypt:      CDATA(encMsg),
		MsgSignature: CDATA(sortSha1(token, timeStamp, nonce, encMsg)),
		TimeStamp:    timeStamp,
		Nonce:        CDATA(nonce),
	}
	return xml.MarshalIndent(re, " ", "  ")
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
