package wechat

import (
	"fmt"
	"log"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 企业号用户列表接口
const (
	WXAPI_GETUSER     = WXAPI_ENT + "user/getuserinfo?access_token=%s&code=%s"
	WXAPI_GETUSERINFO = WXAPI_ENT + "user/get?access_token=%s&userid=%s"
	WXAPI_USERLIST    = WXAPI_ENT + `user/list?access_token=%s&department_id=1&fetch_child=1&status=0`
	WXAPI_USERADD     = WXAPI_ENT + `user/create?access_token=`
	WXAPI_DEPTLIST    = WXAPI_ENT + `department/list?access_token=%s&id=1`
)

// UserOauth 用户鉴权信息
type UserOauth struct {
	WxErr
	UserId   string
	DeviceId string
	OpenId   string
}

// GetUserOauth 通过code鉴权
func GetUserOauth(code string) (userOauth UserOauth, err error) {
	url := fmt.Sprintf(WXAPI_GETUSER, GetAccessToken(), code)
	if err = util.GetJson(url, &userOauth); err != nil {
		return
	}
	if userOauth.ErrCode != 0 {
		err = fmt.Errorf("GetUserId error : errcode=%v , errmsg=%v", userOauth.ErrCode, userOauth.ErrMsg)
	}
	return
}

// UserInfo 用户信息
type UserInfo struct {
	WxErr
	UserId     string
	Name       string
	Department []int
	Position   string
	Mobile     string
	Gender     string
	Email      string
	WeixinId   string
	Avatar     string
	Status     int
	ExtAttr    struct {
		Attrs []struct {
			Name  string
			Value string
		}
	}
}

// GetUserInfo 通过userId获取用户信息
func GetUserInfo(userId string) (userInfo UserInfo, err error) {
	url := fmt.Sprintf(WXAPI_GETUSERINFO, GetAccessToken(), userId)
	if err = util.GetJson(url, &userInfo); err != nil {
		return
	}
	if userInfo.ErrCode != 0 {
		err = fmt.Errorf("GetUserId error : errcode=%v , errmsg=%v", userInfo.ErrCode, userInfo.ErrMsg)
	}
	return
}

// UserList 用户列表
var UserList userList

// UserList 用户列表
type userList struct {
	WxErr
	UserList []UserInfo
}

// UpdateUserList 获取用户列表
func UpdateUserList() (err error) {
	UserList, err = GetUserList()
	if err != nil {
		log.Println("同步通讯录失败:", err)
	}
	return
}

// GetUserList 获取用户列表
func GetUserList() (u userList, err error) {
	url := fmt.Sprintf(WXAPI_USERLIST, GetAccessToken())
	if err = util.GetJson(url, &u); err != nil {
		return
	}
	if u.ErrCode != 0 {
		err = fmt.Errorf("GetUserList error : errcode=%v , errmsg=%v", u.ErrCode, u.ErrMsg)
	}
	return
}

// DeptList 部门列表
var DeptList DepartmentList

// DepartmentList 部门列表
type DepartmentList struct {
	WxErr
	Department []struct {
		Id       int
		Name     string
		ParentId int
		Order    int
	}
}

// UpdateDeptList 更新部门列表
func UpdateDeptList() (err error) {
	DeptList, err = GetDeptList()
	if err != nil {
		log.Println("获取部门列表失败:", err)
	}
	return
}

// GetDeptList 获取部门列表
func GetDeptList() (deptList DepartmentList, err error) {
	url := fmt.Sprintf(WXAPI_DEPTLIST, GetAccessToken())
	if err = util.GetJson(url, &deptList); err != nil {
		return
	}
	if deptList.ErrCode != 0 {
		err = fmt.Errorf("GetDeptList error : errcode=%v , errmsg=%v", deptList.ErrCode, deptList.ErrMsg)
	}
	return
}

// GetDeptName 通过部门id获取部门名称
func GetDeptName(id int) string {
	for _, v := range DeptList.Department {
		if v.Id == id {
			return v.Name
		}
	}
	return ""
}

// GetUser 通过账号获取用户信息
func GetUser(userid string) *UserInfo {
	for _, v := range UserList.UserList {
		if v.UserId == userid {
			return &v
		}
	}
	return nil
}

var toUserReplacer = strings.NewReplacer("|", ",", "，", ",")

// GetToUser 获取acl所包含的所有用户
func GetToUser(acl interface{}) (touser string) {
	s1 := strings.TrimSpace(fmt.Sprint(acl))
	if s1 == "@all" {
		return "@all"
	}
	arr := strings.Split(toUserReplacer.Replace(s1), ",")
	for _, toUser := range arr {
		for _, v := range UserList.UserList {
			if CheckUserAcl(v.UserId, toUser) {
				touser += "|" + v.UserId
			}
		}
	}
	return strings.Trim(touser, "|")
}

// CheckUserAcl 测试权限，对比user的账号，姓名，手机，职位是否包含于acl
func CheckUserAcl(userid, acl string) bool {
	acl = strings.TrimSpace(acl)
	if acl == "" {
		return false
	}
	if acl == "@all" {
		return true
	}
	acl = "," + strings.Replace(acl, "，", ",", -1) + ","
	u := GetUser(userid)
	if u == nil {
		return false
	}
	for _, dv := range u.Department {
		if strings.Contains(acl, ","+GetDeptName(dv)+",") {
			return true
		}
		if strings.Contains(acl, ","+GetDeptName(dv)+"/"+u.Position+",") {
			return true
		}
	}
	for _, pv := range u.ExtAttr.Attrs {
		if strings.Contains(acl, ","+pv.Value+",") {
			return true
		}
	}
	return strings.Contains(acl, ","+u.Name+",") ||
		strings.Contains(acl, ","+u.UserId+",") ||
		strings.Contains(acl, ","+u.Mobile+",") ||
		strings.Contains(acl, ","+u.Position+",")
}
