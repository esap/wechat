// Package wechat TODO：微信支付接口
package wechat

import (
	"fmt"

	"time"

	"github.com/esap/wechat/util"
)

// PayRoot 支付根URL
const (
	PayRoot            = "weixin：//wxpay/bizpayurl?"
	PayUrl             = "weixin：//wxpay/bizpayurl?sign=%s&appid=%s&mch_id=%s&product_id=%sX&time_stamp=%vX&nonce_str=%s"
	PayUnifiedOrderUrl = "https://api.mch.weixin.qq.com/pay/unifiedordefunc"
)

// UnifiedOrderReq 统一下单请求体
type UnifiedOrderReq struct {
	Appid          string `xml:"appid"`
	MchId          string `xml:"mch_id"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	SignType       string `xml:"sign_type"`
	Body           string `xml:"body"`
	Detail         CDATA  `xml:"detail"`
	Attach         string `xml:"attach"`
	OutTradeNo     string `xml:"out_trade_no"`
	FeeType        string `xml:"fee_type"`
	TotalFee       string `xml:"total_fee"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
	TimeStart      string `xml:"time_start"`
	TimeExpire     string `xml:"time_expire"`
	GoodsTag       string `xml:"goods_tag"`
	NotifyUrl      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	ProductId      string `xml:"product_id"`
	LimitPay       string `xml:"limit_pay"`
	Openid         string `xml:"openid"`
	SceneInfo      string `xml:"scene_info"`
}

// UnifiedOrderRet 统一下单返回体
type UnifiedOrderRet struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	// 以下字段在return_code为SUCCESS的时候有返回
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	// 以下字段在return_code 和result_code都为SUCCESS的时候有返回
	TradeType  string `xml:"trade_type"`
	PrepayId   string `xml:"prepay_id"`
	CodeUrl    string `xml:"code_url"`
	TimeExpire string `xml:"time_expire"`
	GoodsTag   string `xml:"goods_tag"`
	NotifyUrl  string `xml:"notify_url"`
	ProductId  string `xml:"product_id"`
	LimitPay   string `xml:"limit_pay"`
	Openid     string `xml:"openid"`
	SceneInfo  string `xml:"scene_info"`
}

// GetUnifedOrderUrl 获取统一下单URL，用于生成付款二维码等
func (s *Server) GetUnifedOrderUrl(desc, tradeNo, fee, ip, callback, tradetype, productid string) string {
	noncestr := util.GetRandomString(16)
	r := &UnifiedOrderReq{
		Appid:          s.AppId,
		MchId:          s.MchId,
		NonceStr:       noncestr,
		Sign:           util.SortMd5(noncestr),
		Body:           desc,
		OutTradeNo:     tradeNo,
		TotalFee:       fee,
		SpbillCreateIp: ip,
		NotifyUrl:      callback,
		TradeType:      tradetype,
		ProductId:      productid,
	}
	ret := new(UnifiedOrderRet)
	err := util.PostXmlPtr(PayUnifiedOrderUrl, r, ret)
	if err != nil {
		Println("GetUnifedOrderUrl err:", err)
		return ""
	}
	return ret.CodeUrl
}

// PayOrderScan 扫码付
func (s *Server) PayOrderScan(mchId, ProductId string) string {
	nonceStr := util.GetRandomString(10)
	timeStamp := time.Now().Unix()
	strA := fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&product_id=%s&time_stamp=%v", s.AppId, mchId, nonceStr, ProductId, timeStamp)
	return PayRoot + strA + "&sign=" + util.SortMd5(strA)
}
