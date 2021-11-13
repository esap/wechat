# WeChat SDK
[![Build Status](https://travis-ci.org/esap/wechat.svg?branch=master)](https://travis-ci.org/esap/wechat)
[![Go Report Card](https://goreportcard.com/badge/github.com/esap/wechat)](https://goreportcard.com/report/github.com/esap/wechat)
[![GoDoc](http://godoc.org/github.com/esap/wechat?status.svg)](http://godoc.org/github.com/esap/wechat)

**微信SDK的golang实现，短小精悍，同时兼容【企业微信/服务号/订阅号/小程序】**

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
	
	cfg := &wechat.WxConfig{
		Token:          "yourToken",
		AppId:          "yourAppID",
		Secret:         "yourSecret",
		EncodingAESKey: "yourEncodingAesKey",
	}

	app := wechat.New(cfg)
	app.SendText("@all", "Hello,World!")
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.VerifyURL(w, r).NewText("客服消息1").Send().NewText("客服消息2").Send().NewText("查询OK").Reply()
	})
	
	http.ListenAndServe(":9090", nil)
}

```
## 配置方式

```go
	// 创建公众号实例(服务号/订阅号/小程序) 不带aesKey则为明文模式
	cfg := &wechat.WxConfig{
		Token:          "yourToken",
		AppId:          "yourAppID",
		Secret:         "yourSecret",
	}

	// 创建公众号实例(服务号/订阅号/小程序)
	cfg := &wechat.WxConfig{
		Token:          "yourToken",
		AppId:          "yourAppID",
		Secret:         "yourSecret",
		EncodingAESKey: "yourEncodingAesKey",
	}

	// 创建企业微信实例
	cfg := &wechat.WxConfig{
		Token:          "yourToken",
		AppId:          "yourCorpID",
		AgentId:        "yourAgentId",
		Secret:         "yourSecret",
		EncodingAESKey: "yourEncodingAesKey",
		AppType:        1,
	}
```

## 主动推送消息

用户关注后，企业微信可以主动推送消息，服务号需要用户48小时内进入过。

```go
	app.SendText(to, msg)
	app.SendImage(to, mediaId)
	app.SendVoice(to, mediaId)
	app.SendFile(to, mediaId)
	app.SendVideo(to, mediaId, title, desc)
	app.SendTextcard(to, title, desc, url)
	app.SendMusic(to, mediaId, title, desc, musicUrl, qhMusicUrl)
	app.SendNews(to, arts...)
	app.SendMpNews(to, arts...)
	app.SendMpNewsId(to, mediaId)
	app.SendMarkDown(to, content)
```

## 消息回调

* 通常将`app.VerifyURL(http.ResponseWriter, *http.Request)`嵌入http handler

该函数返回`*wechat.Context`基本对象，其中的Msg为用户消息：

```go
	// 混合用户消息，业务判断的主体
	WxMsg struct {
		XMLName      xml.Name `xml:"xml"`
		ToUserName  
		FromUserName
		CreateTime  64
		MsgId       64
		MsgType     
		Content       // text
		AgentID          // corp
		PicUrl        // image
		MediaId       // image/voice/video/shortvideo
		Format        // voice
		Recognition   // voice
		ThumbMediaId  // video
		LocationX    float32 `xml:"Latitude"`  // location
		LocationY    float32 `xml:"Longitude"` // location
		Precision    float32 // LOCATION
		Scale            // location
		Label         // location
		Title         // link
		Description   // link
		Url           // link
		Event         // event
		EventKey      // event
		SessionFrom   // event|user_enter_tempsession
		Ticket      
		FileKey     
		FileMd5     
		FileTotalLen

		ScanCodeInfo struct {
			ScanType  
			ScanResult
		}
	}

```

* 如果使用其他web框架，例如echo/gin/beego等，则把VerifyURL()放入controller或handler

```go
// echo示例 公众号回调接口
func wxApiPost(c echo.Context) error {
	ctx := app.VerifyURL(c.Response().Writer, c.Request())
	
	// TODO: 这里是其他业务操作

	return nil
}
```

### 回调回复消息

回调回复消息有两种方式：

* 被动回复，采用XML格式编码返回(Reply)；

* 客服消息，采用json格式编码返回(Send)；

* 两种方式都可先调用`*wechat.Context`对象的New方法创建消息，然后调用Reply()或Send()。

* 支持链式调用，但Reply()只有第一次有效。

```go
	ctx.NewText("正在查询中...").Reply()
	ctx.NewText("客服消息1").Send().NewText("客服消息2").Send()
```

* 被动回复可直接调用Reply()，表示已收到，然后调用客服消息。

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

####  模板消息

[相关issue](https://github.com/esap/wechat/issues/20#issue-451068915)

```go
	tlpdata := map[string]struct {
		Value `json:"value"`
		Color `json:"color"`
	}{
		"first": {Value: "我是渣渣涛", Color: "#173177"},
		"keyword1": {Value: "这是一个你从没有玩过的全新游戏", Color: "#173177"},
		"keyword2": {Value: "只要你跟着我一起试玩一下", Color: "#173177"},
		"keyword3": {Value: "你就会爱上这款游戏", Color: "#4B1515"},
		"remark":   {Value: "是兄弟就来砍我", Color: "#071D42"},
	}
	ctx.SendTemplate(
		ctx.Msg.FromUserName,
		"tempid", // 模板ID
		c.Request.Host, // 跳转url
		ctx.AppId, // 跳转小程序，比url优先
		"", // 小程序页面
		tlpdata,
	)
```

## License

MIT
