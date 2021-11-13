package wechat

import (
	"fmt"

	"github.com/esap/wechat/util"
)

// MPUserGetList 公众号用户接口
const (
	MPUserGetList  = WXAPI + "user/get?access_token=%s&next_openid=%s"
	MPUserBatchGet = WXAPI + "user/info/batchget?access_token="
	MPUserInfo     = WXAPI + "user/info?access_token=%s&openid=%v&lang=%v"
)

type (
	// MpUserInfoList 公众号用户信息列表
	MpUserInfoList struct {
		WxErr
		MpUserInfoList []MpUserInfo `json:"user_info_list"`
	}

	// MpUserInfo 公众号用户信息
	MpUserInfo struct {
		Subscribe     int
		OpenId        string
		NickName      string
		Sex           int
		Language      string
		City          string
		Province      string
		Country       string
		HeadImgUrl    string
		SubscribeTime int `json:"subscribe_time"`
		UnionId       string
		Remark        string
		GroupId       int
		TagIdList     []int `json:"tagid_list"`
	}

	// MpUser 服务号用户
	MpUser struct {
		WxErr
		Total int
		Count int
		Data  struct {
			OpenId []string
		}
		NextOpenId string
	}

	// MpUserListReq 公众号用户请求
	MpUserListReq struct {
		UserList interface{} `json:"user_list"`
	}
)

// BatchGetAll 获取所有公众号用户
func (s *Server) BatchGetAll() (ui []MpUserInfo, err error) {
	var ul []string
	ul, err = s.GetAllMpUserList()
	if err != nil {
		return
	}
	leng := len(ul)
	if leng <= 100 {
		return s.BatchGet(ul)
	}
	for i := 0; i < leng/100+1; i++ {
		end := (i + 1) * 100
		if end > leng {
			end = leng
		}

		ui2, err2 := s.BatchGet(ul[i*100 : end])
		if err != nil {
			err = err2
			return
		}
		ui = append(ui, ui2...)
	}
	return
}

// BatchGet 批量获取公众号用户信息
func (s *Server) BatchGet(ul []string) (ui []MpUserInfo, err error) {
	m := make([]map[string]interface{}, len(ul))

	for k, v := range ul {
		m[k] = make(map[string]interface{})
		m[k]["openid"] = v
	}
	ml := new(MpUserInfoList)
	err = util.PostJsonPtr(MPUserBatchGet+s.GetAccessToken(), MpUserListReq{m}, ml)
	return ml.MpUserInfoList, ml.Error()
}

// GetAllMpUserList 获取所有用户ID
func (s *Server) GetAllMpUserList() (ul []string, err error) {
	ul = make([]string, 0)
	mul, err := s.GetMpUserList()
	if err != nil {
		return
	}
	if mul.Error() == nil {
		ul = append(ul, mul.Data.OpenId...)
	}
	for mul.Count == 10000 {
		mul, err = s.GetMpUserList(mul.NextOpenId)
		if err != nil {
			return
		}
		if mul.Error() == nil {
			ul = append(ul, mul.Data.OpenId...)
		}
	}
	return
}

// GetMpUserList 获取用户信息，根据openid
func (s *Server) GetMpUserList(openid ...string) (ul *MpUser, err error) {
	if len(openid) == 0 {
		openid = append(openid, "")
	}
	mpuser := new(MpUser)
	url := fmt.Sprintf(MPUserGetList, s.GetAccessToken(), openid[0])
	if err = util.GetJson(url, &mpuser); err != nil {
		return
	}
	return mpuser, mpuser.Error()
}

// GetMpUserInfo 获取用户详情
func (s *Server) GetMpUserInfo(openid string, lang ...string) (user *MpUserInfo, err error) {
	if len(lang) == 0 {
		lang = append(lang, "zh_CN")
	}
	user = new(MpUserInfo)
	url := fmt.Sprintf(MPUserInfo, s.GetAccessToken(), openid, lang[0])
	if err = util.GetJson(url, &user); err != nil {
		return
	}
	return
}
