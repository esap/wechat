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
		wechat.Set("yourToken", "yourAppID", "yourSecret", "yourEncodingAesKey")
		http.HandleFunc("/", WxHandler)
		http.ListenAndServe(":9090", nil)
	}

	func WxHandler(w http.ResponseWriter, r *http.Request) {
		wechat.VerifyURL(w, r).NewText("客服消息1").Send().NewText("客服消息2").Send().NewText("查询结果...").Reply()
	}

More info: https://github.com/esap/wechat

*/
package wechat
