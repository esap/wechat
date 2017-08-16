package wechat

import (
	"fmt"
	"log"

	"github.com/esap/wechat/util"
)

// WXAPI 企业号标签接口
const (
	WXAPI_TagList     = WXAPI_ENT + `tag/list?access_token=%s`
	WXAPI_TagUsers    = WXAPI_ENT + `tag/get?access_token=`
	WXAPI_AddTagUsers = WXAPI_ENT + `tag/addtagusers?access_token=`
	WXAPI_DelTagUsers = WXAPI_ENT + `tag/deltagusers?access_token=`
	WXAPI_TagAdd      = WXAPI_ENT + `tag/create?access_token=`
	WXAPI_TagUpdate   = WXAPI_ENT + `tag/update?access_token=`
	WXAPI_TagDel      = WXAPI_ENT + `tag/delete?access_token=%s&id=%d`
)

// DeptList 标签列表
var Tags TagList

type (
	// TagList 标签列表
	TagList struct {
		WxErr
		Taglist []Tag
	}

	// Tag 标签
	Tag struct {
		TagId   int    `json:"tagid"`
		TagName string `json:"tagname"`
	}
	// TagUsers 标签成员
	TagUsers struct {
		WxErr
		TagName   string
		UserList  []UserInfo
		PartyList []int
	}
)

// SyncTagList 更新标签列表
func (s *Server) SyncTagList() (err error) {
	s.TagList, err = s.GetTagList()
	if err != nil {
		log.Println("获取标签列表失败:", err)
	}
	return
}

// GetTagList 获取标签列表
func (s *Server) GetTagList() (l TagList, err error) {
	url := fmt.Sprintf(WXAPI_TagList, s.GetAccessToken())
	if err = util.GetJson(url, &l); err != nil {
		return
	}
	if l.ErrCode != 0 {
		err = fmt.Errorf("GetTagList error : errcode=%v , errmsg=%v", l.ErrCode, l.ErrMsg)
	}
	return
}

// TagAdd 获取标签列表
func (s *Server) TagAdd(Tag *Tag) (err error) {
	return s.doUpdate(WXAPI_TagAdd, Tag)
}

// TagUpdate 获取标签列表
func (s *Server) TagUpdate(Tag *Tag) (err error) {
	return s.doUpdate(WXAPI_TagUpdate, Tag)
}

// TagDelete 删除用户
func (s *Server) TagDelete(TagId string) (err error) {
	e := new(WxErr)
	if err = util.GetJson(WXAPI_TagDel+s.GetAccessToken()+"&tagid="+TagId, e); err != nil {
		return
	}
	return e.Error()
}

// GetTagUsers 获取标签下的成员
func (s *Server) GetTagUsers(id string) (tm *TagUsers, err error) {
	err = util.GetJson(WXAPI_TagUsers+s.GetAccessToken()+"&tagid="+id, tm)
	return
}

// AddTagUsers 添加标签成员
func (s *Server) AddTagUsers(id string, userlist, partylist []string) error {
	m := map[string]interface{}{"tagid": id, "userlist": userlist, "partylist": partylist}
	return s.doUpdate(WXAPI_AddTagUsers, m)
}

// DelTagUsers 删除标签成员
func (s *Server) DelTagUsers(id string, userlist, partylist []string) error {
	m := map[string]interface{}{"tagid": id, "userlist": userlist, "partylist": partylist}
	return s.doUpdate(WXAPI_DelTagUsers, m)
}

// GetTagName 通过标签id获取标签名称
func (s *Server) GetTagName(id int) string {
	for _, v := range s.TagList.Taglist {
		if v.TagId == id {
			return v.TagName
		}
	}
	return ""
}
