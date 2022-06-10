package wechat

import (
	"fmt"
	"log"
	"strings"

	"github.com/esap/wechat/util"
)

// CorpAPIDeptList 企业微信部门列表接口
const (
	CorpAPIDeptList   = CorpAPI + `department/list?access_token=%s`
	CorpAPIDeptAdd    = CorpAPI + `department/create?access_token=`
	CorpAPIDeptUpdate = CorpAPI + `department/update?access_token=`
	CorpAPIDeptDel    = CorpAPI + `department/delete?access_token=`
)

type (
	// DeptList 部门列表
	DeptList struct {
		WxErr
		Department []Department
	}

	// Department 部门
	Department struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		ParentId int    `json:"parentid"`
		Order1   int64  `json:"order"`
	}
)

// SyncDeptList 更新部门列表
func (s *Server) SyncDeptList() (err error) {
	s.DeptList, err = s.GetDeptList()
	if err != nil {
		log.Printf("[%v::%v]获取部门列表失败:%v", s.AppId, s.AgentId, err)
	}
	return
}

// GetDeptList 获取部门列表
func (s *Server) GetDeptList() (dl DeptList, err error) {
	url := fmt.Sprintf(CorpAPIDeptList, s.GetUserAccessToken())
	if err = util.GetJson(url, &dl); err != nil {
		return
	}
	err = dl.Error()
	return
}

// GetDeptIdList 获取部门id列表
func (s *Server) GetDeptIdList() (deptIdlist []int) {
	deptIdlist = make([]int, 0)
	s.SyncDeptList()
	for _, v := range s.DeptList.Department {
		deptIdlist = append(deptIdlist, v.Id)
	}
	return
}

// DeptAdd 获取部门列表
func (s *Server) DeptAdd(dept *Department) (err error) {
	return s.doUpdate(CorpAPIDeptAdd, dept)
}

// DeptUpdate 获取部门列表
func (s *Server) DeptUpdate(dept *Department) (err error) {
	return s.doUpdate(CorpAPIDeptUpdate, dept)
}

// DeptDelete 删除部门
func (s *Server) DeptDelete(Id int) (err error) {
	e := new(WxErr)
	if err = util.GetJson(CorpAPIDeptDel+s.GetUserAccessToken()+"&id="+fmt.Sprint(Id), e); err != nil {
		return
	}
	return e.Error()
}

// GetDeptName 通过部门id获取部门名称
func (s *Server) GetDeptName(id int) string {
	for _, v := range s.DeptList.Department {
		if v.Id == id {
			return v.Name
		}
	}
	return ""
}

// GetToParty 获取acl所包含的所有部门ID，结果形式：tagId1|tagId2|tagId3...
func (s *Server) GetToParty(acl interface{}) string {
	s1 := strings.TrimSpace(acl.(string))
	arr := strings.Split(toUserReplacer.Replace(s1), "|")
	for k, totag := range arr {
		for _, v := range s.DeptList.Department {
			if v.Name == totag {
				arr[k] = fmt.Sprint(v.Id)
			}
		}
	}
	return strings.Join(arr, "|")
}

// CheckDeptAcl 测试权限，对比user是否包含于acl
func (s *Server) CheckDeptAcl(userid, acl string) bool {
	acl = strings.TrimSpace(acl)
	if acl == "" {
		return false
	}
	u := s.GetUser(userid)
	if u == nil {
		return false
	}
	acl = "|" + toUserReplacer.Replace(acl) + "|"
	for _, id := range u.Department {
		if strings.Contains(acl, "|"+s.GetDeptName(id)+"|") {
			return true
		}
		if strings.Contains(acl, "|"+fmt.Sprint(id)+"|") {
			return true
		}
	}

	return false
}
