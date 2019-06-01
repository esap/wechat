package wechat

import (
	"errors"
	"fmt"

	"github.com/esap/wechat/util"
)

// MP_Template 公众号模板消息接口
const (
	MP_AddTemplate     = WXAPI + "template/api_add_template?access_token="
	MP_GetAllTemplate  = WXAPI + "template/get_all_private_template?access_token="
	MP_DelTemplate     = WXAPI + "template/del_private_template?access_token="
	MP_SendTemplateMsg = WXAPI + "message/template/send?access_token="
)

// MpTemplate 模板信息
type MpTemplate struct {
	TemplateId      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"template_id"`
	Example         string `json:"example"`
}

// AddTemplate 获取模板
func (s *Server) AddTemplate(IdShort string) (id string, err error) {
	form := map[string]interface{}{"template_id_short": IdShort}

	ret := make(map[string]interface{})
	err = util.PostJsonPtr(MP_AddTemplate+s.GetAccessToken(), form, ret)
	if err != nil {
		return
	}

	if fmt.Sprint(ret["errcode"]) != "0" {
		return "", errors.New(fmt.Sprint(ret["errcode"]))
	}

	return ret["template_id"].(string), nil
}

// DelTemplate 删除模板
func (s *Server) DelTemplate(id string) (err error) {
	form := map[string]interface{}{"template_id": id}

	ret := make(map[string]interface{})
	err = util.PostJsonPtr(MP_DelTemplate+s.GetAccessToken(), form, ret)
	if err != nil {
		return
	}

	if fmt.Sprint(ret["errcode"]) != "0" {
		return errors.New(fmt.Sprint(ret["errcode"]))
	}

	return
}

// GetAllTemplate 获取模板
func (s *Server) GetAllTemplate() (templist []MpTemplate, err error) {
	ret := make(map[string]interface{})
	err = util.GetJson(MP_GetAllTemplate+s.GetAccessToken(), ret)
	if err != nil {
		return
	}

	if fmt.Sprint(ret["errcode"]) != "0" {
		return nil, errors.New(fmt.Sprint(ret["errcode"]))
	}

	return ret["template_id"].([]MpTemplate), nil
}

// SendTemplate 发送模板消息，data通常是map[string]struct{value string,color string}
func (s *Server) SendTemplate(to, id, url, appid, pagepath string, data interface{}) (msgid float64, err error) {

	form := map[string]interface{}{
		"touser":      to,
		"template_id": id,
		"data":        data,
	}
	if pagepath != "" {
		form["miniprogram"] = map[string]string{
			"appid":    appid,
			"pagepath": pagepath,
		}
	} else if url != "" {
		form["url"] = url
	}
	ret := make(map[string]interface{})
	err = util.PostJsonPtr(MP_SendTemplateMsg+s.GetAccessToken(), form, &ret)
	if err != nil {
		return
	}

	if fmt.Sprint(ret["errcode"]) != "0" {
		return 0, (errors.New(fmt.Sprint(ret["errcode"])))
	}

	return ret["msgid"].(float64), nil
}
