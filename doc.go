/*
Package wechat provide wechat-sdk for go

5行代码，开启微信API示例:

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.VerifyURL(w, r).NewText("客服消息1").Send().NewText("客服消息2").Send().NewText("查询OK").Reply()
	})

	http.ListenAndServe(":9090", nil)
}

More info: https://github.com/esap/wechat

*/
package wechat
