package wechat

import (
	"fmt"

	"github.com/esap/wechat/util"
)

// WXAPI_ENT 企业号菜单接口
const (
	WXAPI_GetMenu = `menu/get?access_token=%s&agentid=%d`
	WXAPI_AddMenu = `menu/create?access_token=%s&agentid=%d`
	WXAPI_DelMenu = `menu/delete?access_token=%s&agentid=%d`
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

// GetMenu 获取应用菜单
func (s *Server) GetMenu() (m *Menu, err error) {
	m = new(Menu)
	url := fmt.Sprintf(s.RootUrl+WXAPI_GetMenu, s.GetAccessToken(), s.AgentId)
	if err = util.GetJson(url, m); err != nil {
		return
	}
	err = m.Error()
	return
}

// AddMenu 创建应用菜单
func (s *Server) AddMenu(m *Menu) (err error) {
	e := new(WxErr)
	url := fmt.Sprintf(s.RootUrl+WXAPI_AddMenu, s.GetAccessToken(), s.AgentId)
	if err = util.PostJsonPtr(url, m, e); err != nil {
		return
	}
	return e.Error()
}

// DelMenu 删除应用菜单
func (s *Server) DelMenu() (err error) {
	e := new(WxErr)
	url := fmt.Sprintf(s.RootUrl+WXAPI_DelMenu, s.GetAccessToken(), s.AgentId)
	if err = util.GetJson(url, e); err != nil {
		return
	}
	return e.Error()
}
