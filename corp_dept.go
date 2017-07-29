package wechat

import (
	"fmt"
	"log"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 企业号部门列表接口
const (
	WXAPI_DeptList   = WXAPI_ENT + `department/list?access_token=%s&id=1`
	WXAPI_DeptAdd    = WXAPI_ENT + `department/create?access_token=`
	WXAPI_DeptUpdate = WXAPI_ENT + `department/update?access_token=`
	WXAPI_DeptDel    = WXAPI_ENT + `department/delete?access_token=%s&id=%d`
)

// DeptList 部门列表
//var Depts DeptList

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
		log.Println("获取部门列表失败:", err)
	}
	return
}

// GetDeptList 获取部门列表
func (s *Server) GetDeptList() (dl DeptList, err error) {
	url := fmt.Sprintf(WXAPI_DeptList, s.GetAccessToken())
	if err = util.GetJson(url, &dl); err != nil {
		return
	}
	if dl.ErrCode != 0 {
		err = fmt.Errorf("GetDeptList error : errcode=%v , errmsg=%v", dl.ErrCode, dl.ErrMsg)
	}
	return
}

// GetDeptIdList 获取部门id列表
func (s *Server) GetDeptIdList() (deptIdlist []int) {
	deptIdlist = make([]int, 0)
	for _, v := range s.DeptList.Department {
		deptIdlist = append(deptIdlist, v.Id)
	}
	return
}

// DeptAdd 获取部门列表
func (s *Server) DeptAdd(dept *Department) (err error) {
	return s.doUpdate(WXAPI_DeptAdd, dept)
}

// DeptUpdate 获取部门列表
func (s *Server) DeptUpdate(dept *Department) (err error) {
	return s.doUpdate(WXAPI_DeptUpdate, dept)
}

// DeptDelete 删除部门
func (s *Server) DeptDelete(Id int) (err error) {
	e := new(WxErr)
	if err = util.GetJson(WXAPI_DeptDel+s.GetAccessToken()+"&id="+fmt.Sprint(Id), e); err != nil {
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
