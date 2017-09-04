package wechat

import (
	"fmt"
	"log"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 企业号标签接口
const (
	WXAPI_TagList     = WXAPI_ENT + `tag/list?access_token=`
	WXAPI_TagUsers    = WXAPI_ENT + `tag/get?access_token=`
	WXAPI_AddTagUsers = WXAPI_ENT + `tag/addtagusers?access_token=`
	WXAPI_DelTagUsers = WXAPI_ENT + `tag/deltagusers?access_token=`
	WXAPI_TagAdd      = WXAPI_ENT + `tag/create?access_token=`
	WXAPI_TagUpdate   = WXAPI_ENT + `tag/update?access_token=`
	WXAPI_TagDel      = WXAPI_ENT + `tag/delete?access_token=`
)

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
		TagId     int `json:"tagid"`
		TagName   string
		UserList  []UserInfo
		PartyList []int
	}

	// TagUserBody 标签成员（请求body格式）
	TagUserBody struct {
		TagId     int      `json:"tagid"`
		UserList  []string `json:"userlist"`
		PartyList []int    `json:"partylist"`
	}

	// TagErr
	TagErr struct {
		WxErr
		InvalidList  string
		InvalidParty []int
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
	url := WXAPI_TagList + s.GetAccessToken()
	if err = util.GetJson(url, &l); err != nil {
		return
	}
	err = l.Error()
	return
}

// GetTagIdList 获取标签id列表
func (s *Server) GetTagIdList() (tagIdlist []int) {
	tagIdlist = make([]int, 0)
	for _, v := range s.TagList.Taglist {
		tagIdlist = append(tagIdlist, v.TagId)
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
func (s *Server) TagDelete(TagId int) (err error) {
	e := new(WxErr)
	if err = util.GetJson(WXAPI_TagDel+s.GetAccessToken()+"&tagid="+fmt.Sprint(TagId), e); err != nil {
		return
	}
	return e.Error()
}

// GetTagUsers 获取标签下的成员
func (s *Server) GetTagUsers(id int) (tu *TagUsers, err error) {
	tu = new(TagUsers)
	err = util.GetJson(WXAPI_TagUsers+s.GetAccessToken()+"&tagid="+fmt.Sprint(id), tu)
	return
}

// AddTagUsers 添加标签成员
func (s *Server) AddTagUsers(id int, userlist []string, partylist []int) error {
	leng := len(userlist)
	e := new(TagErr)
	for i := 0; i < leng/1000+1; i++ {
		end := (i + 1) * 1000
		if end > leng {
			end = leng
		}
		b := TagUserBody{TagId: id, UserList: userlist[i*1000 : end], PartyList: partylist}
		url := WXAPI_AddTagUsers + s.GetAccessToken()
		if err := util.PostJsonPtr(url, b, e); err != nil {
			return err
		}
	}
	return e.Error()
}

// DelTagUsers 删除标签成员
func (s *Server) DelTagUsers(id string, body TagUserBody) error {
	return s.doUpdate(WXAPI_DelTagUsers, body)
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

// GetTagId 通过标签名称获取标签名称
func (s *Server) GetTagId(name string) int {
	for _, v := range s.TagList.Taglist {
		if fmt.Sprint(v.TagId) == name || v.TagName == name {
			return v.TagId
		}
	}
	return 0
}

// GetToTag 获取acl所包含的所有标签ID，结果形式：tagId1|tagId2|tagId3...
func (s *Server) GetToTag(acl interface{}) string {
	s1 := strings.TrimSpace(acl.(string))
	arr := strings.Split(toUserReplacer.Replace(s1), "|")
	for k, totag := range arr {
		for _, v := range s.TagList.Taglist {
			if v.TagName == totag {
				arr[k] = fmt.Sprint(v.TagId)
			}
		}
	}
	return strings.Join(arr, "|")
}

// CheckTagrAcl 测试权限，对比user是否包含于acl
func (s *Server) CheckTagAcl(userid, acl string) bool {
	acl = strings.TrimSpace(acl)
	if acl == "" {
		return false
	}
	arr := strings.Split(toUserReplacer.Replace(acl), "|")
	for _, idOrName := range arr {
		tu, err := s.GetTagUsers(s.GetTagId(idOrName))
		if err != nil {
			continue
		}
		for _, u := range tu.UserList {
			if u.UserId == userid {
				return true
			}
		}
	}
	return false
}
