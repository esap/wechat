package wechat

import (
	"fmt"

	"github.com/esap/wechat/util"
)

// WXAPI_ENT 企业号菜单接口
const (
	WXAPI_GetCorpMenu = WXAPI_ENT + `menu/get?access_token=%s&agentid=%d`
	WXAPI_AddCorpMenu = WXAPI_ENT + `menu/create?access_token=%s&agentid=%d`
	WXAPI_DelCorpMenu = WXAPI_ENT + `menu/delete?access_token=%s&agentid=%d`
)

type (
	// Menu 菜单
	Menu struct {
		WxErr
		Button []struct {
			Name      string `json:"name"`
			Type      string `json:"type"`
			Key       string `json:"key"`
			Url       string `json:"url"`
			SubButton []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				Key  string `json:"key"`
				Url  string `json:"url"`
			} `json:"sub_button"`
		} `json:"button"`
	}
)

// GetCorpMenu 获取企业应用菜单
func (s *Server) GetCorpMenu() (m *Menu, err error) {
	url := fmt.Sprintf(WXAPI_GetCorpMenu, s.GetAccessToken(), s.AgentId)
	if err = util.GetJson(url, m); err != nil {
		return
	}
	err = m.Error()
	return
}

// AddCorpMenu 创建企业应用菜单
func (s *Server) AddCorpMenu(m *Menu) (err error) {
	e := new(WxErr)
	url := fmt.Sprintf(WXAPI_AddCorpMenu, s.GetAccessToken(), s.AgentId)
	if err = util.PostJsonPtr(url, m, e); err != nil {
		return
	}
	return e.Error()
}

// DelCorpMenu 删除企业应用菜单
func (s *Server) DelCorpMenu() (err error) {
	e := new(WxErr)
	url := fmt.Sprintf(WXAPI_DelCorpMenu, s.GetAccessToken(), s.AgentId)
	if err = util.GetJson(url, e); err != nil {
		return
	}
	return e.Error()
}
