package wechat

import (
	"errors"
	"fmt"

	"github.com/esap/wechat/util"
)

// MPTemplateGetAll 服务号模板消息接口
const (
	MPTemplateGetAll  = WXAPI + "template/get_all_private_template?access_token="
	MPTemplateAdd     = WXAPI + "template/api_add_template?access_token="
	MPTemplateDel     = WXAPI + "template/del_private_template?access_token="
	MPTemplateSendMsg = WXAPI + "message/template/send?access_token="
)

// MpTemplate 模板信息
type MpTemplate struct {
	TemplateId      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

// AddTemplate 获取模板
func (s *Server) AddTemplate(IdShort string) (id string, err error) {
	form := map[string]interface{}{"template_id_short": IdShort}

	ret := make(map[string]interface{})
	err = util.PostJsonPtr(MPTemplateAdd+s.GetAccessToken(), form, ret)
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
	err = util.PostJsonPtr(MPTemplateDel+s.GetAccessToken(), form, ret)
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
	err = util.GetJson(MPTemplateGetAll+s.GetAccessToken(), ret)
	if err != nil {
		return
	}

	if fmt.Sprint(ret["errcode"]) != "0" {
		return nil, errors.New(fmt.Sprint(ret["errcode"]))
	}

	return ret["template_id"].([]MpTemplate), nil
}

// SendTemplate 发送模板消息，data通常是map[string]struct{value string,color string}
func (s *Server) SendTemplate(to, id, url, appid, pagepath string, data interface{}) *WxErr {

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
	ret := new(WxErr)
	err := util.PostJsonPtr(MPTemplateSendMsg+s.GetAccessToken(), form, &ret)
	if err != nil {
		return &WxErr{ErrCode: -1, ErrMsg: err.Error()}
	}

	return ret
}
