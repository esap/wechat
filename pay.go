package wechat

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	PayRoot = "weixin：//wxpay/bizpayurl?"
	PayUrl  = "weixin：//wxpay/bizpayurl?sign=%s&appid=%s&mch_id=%s&product_id=%sX&time_stamp=%vX&nonce_str=%s"
)

func (s Server) PayOrderScan(mchId, ProductId string) string {
	nonceStr := GetRandomString(10)
	timeStamp := time.Now().Unix()
	strA := fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&product_id=%s&time_stamp=%v", s.AppId, mchId, nonceStr, ProductId, timeStamp)
	return PayRoot + strA + "&sign=" + sortMd5(strA)
}

// sortMd5 排序并md5，主要用于计算sign
func sortMd5(s ...string) string {
	sort.Strings(s)
	h := md5.New()
	h.Write([]byte(strings.Join(s, "")))
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

// GetRandomString 获得随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
