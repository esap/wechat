# WeChat SDK
[![Build Status](https://travis-ci.org/esap/wechat.svg?branch=master)](https://travis-ci.org/esap/wechat)
[![Go Report Card](https://goreportcard.com/badge/github.com/esap/wechat)](https://goreportcard.com/report/github.com/esap/wechat)
[![GoDoc](http://godoc.org/github.com/esap/wechat?status.svg)](http://godoc.org/github.com/esap/wechat)

**微信SDK的golang实现，短小精悍，同时兼容【企业号/服务号/订阅号/小程序】**

## 快速开始

5行代码，链式消息，快速开启微信API示例:

```go
package main

import (
	"net/http"
	"github.com/esap/wechat" // 微信SDK包
)

func main() {
	wechat.Debug = true
	wechat.Set("yourToken", "yourAppID", "yourSecret", "yourAesKey")
	http.HandleFunc("/", WxHandler)
	http.ListenAndServe(":9090", nil)
}

func WxHandler(w http.ResponseWriter, r *http.Request) {
	wechat.VerifyURL(w, r).NewText("这是客服消息").Send().NewText("这是被动回复").Reply()
}

```
## 配置方式

```go
	// 不带aesKey则为明文模式
	wechat.Set("token", "appId", "secrect")

	// 带aesKey则为密文模式
	wechat.Set("token", "appId", "secrect", "aesKey")

	// 企业号自动配置为密文模式
	wechat.SetEnt("token", "appId", "secrect", "aesKey")
```

## 消息管理

通常将`wechat.VerifyURL(http.ResponseWriter, *http.Request)`嵌入http handler

该函数返回`*wechat.Context`基本对象，其中的Msg为用户消息：

```go
// 混合用户消息，业务判断的主体
type WxMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string // text
	MsgId        int64
	AgentID      int     // corp
	PicUrl       string  // image
	MediaId      string  // image/voice/video/shortvideo
	Format       string  // voice
	Recognition  string  // voice
	ThumbMediaId string  // video
	Location_X   float32 `xml:"Latitude"`  // location
	Location_Y   float32 `xml:"Longitude"` // location
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
}

```

### 回复消息

回复消息有两种方式：

* 被动回复，采用XML格式编码返回(Reply)；

* 客服消息，采用json格式编码返回(Send)；

* 两种方式都可先调用`*wechat.Context`对象的New方法创建消息，然后调用Reply()或Send()。

* 支持链式调用，但Reply()只有第一次有效。

```go
	ctx.NewText("正在查询中...").Reply()
	ctx.NewText("客服消息1").Send().NewText("客服消息2").Send()
```

* 被动回复可直接调用ReplySuccess()，表示已收到，然后调用客服消息。

####  文本消息

```go
	ctx.NewText("content")
```

####  图片/语言/文件消息

```go
	// mediaID 可通过素材管理-上上传多媒体文件获得
	ctx.NewImage("mediaID")
	ctx.NewVoice("mediaID")
	
	// 仅企业号支持
	ctx.NewFile("mediaID")
```

####  视频消息

```go
	ctx.NewVideo("mediaID", "title", "description")
```

####  音乐消息

```go
	ctx.NewMusic("thumbMediaID","title", "description", "musicURL", "hqMusicURL")
```

####  图文消息

```go
	// 先创建三个文章
	art1 := wechat.NewArticle("拥抱AI，享受工作",
		"来自村长的ESAP系统最新技术分享",
		"http://ylin.wang/img/esap18-1.png",
		"http://ylin.wang/2017/07/13/esap18/")
	art2 := wechat.NewArticle("用企业微信代替pda实现扫描入库",
		"来自村长的ESAP系统最新技术分享",
		"http://ylin.wang/img/esap17-2.png",
		"http://ylin.wang/2017/06/23/esap17/")
	art3 := wechat.NewArticle("大道至简的哲学",
		"来自村长的工作日志",
		"http://ylin.wang/img/golang.jpg",
		"http://ylin.wang/2017/01/29/log7/")
	// 打包成新闻
	ctx.NewNews(art1, art2, art3)
```

## License

MIT
