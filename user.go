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
	WXAPI_USERUPDATE  = WXAPI_ENT + `user/update?access_token=`
	WXAPI_USERDEL     = WXAPI_ENT + `user/delete?access_token=`
	WXAPI_DEPTLIST    = WXAPI_ENT + `department/list?access_token=%s&id=1`
	WXAPI_DEPTADD     = WXAPI_ENT + `department/create?access_token=`
	WXAPI_DEPTUPDATE  = WXAPI_ENT + `department/update?access_token=`
	WXAPI_DEPTDEL     = WXAPI_ENT + `department/delete?access_token=%s&id=%d`
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
	WxErr      `json:"-"`
	UserId     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	Dept       int    `json:"dept"`
	Position   string `json:"position,omitempty"`
	Mobile     string `json:"mobile"`
	Gender     string `json:"gender,omitempty"`
	Email      string `json:"email,omitempty"`
	WeixinId   string `json:"-"`
	Avatar     string `json:"avatar_mediaid,omitempty"`
	Status     int    `json:"-"`
	ExtAttr    struct {
		Attrs []struct {
			Name  string
			Value string
		}
	} `json:"-"`
}

// UserAdd 添加用户
func UserAdd(user *UserInfo) (err error) {
	return doUpdate(WXAPI_USERADD, user)
}

// UserUpdate 添加用户
func UserUpdate(user *UserInfo) (err error) {
	return doUpdate(WXAPI_USERUPDATE, user)
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

// GetUser 通过账号获取用户信息
func GetUser(userid string) *UserInfo {
	for _, v := range UserList.UserList {
		if v.UserId == userid {
			return &v
		}
	}
	return nil
}

// GetUserName 通过账号获取用户信息
func GetUserName(userid string) string {
	for _, v := range UserList.UserList {
		if v.UserId == userid {
			return v.Name
		}
	}
	return ""
}

// UserList 用户列表
var UserList userList

// UserList 用户列表
type userList struct {
	WxErr
	UserList []UserInfo
}

// SyncUserList 获取用户列表
func SyncUserList() (err error) {
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

// GetUserNameList 获取用户列表
func GetUserNameList() (userlist []string) {
	userlist = make([]string, 0)
	for _, v := range UserList.UserList {
		userlist = append(userlist, v.UserId)
	}
	return
}

// DeptList 部门列表
var DeptList DepartmentList

type (
	// DepartmentList 部门列表
	DepartmentList struct {
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
func SyncDeptList() (err error) {
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

// DeptAdd 获取部门列表
func DeptAdd(dept *Department) (err error) {
	return doUpdate(WXAPI_DEPTADD, dept)
}

// DeptUpdate 获取部门列表
func DeptUpdate(dept *Department) (err error) {
	return doUpdate(WXAPI_DEPTUPDATE, dept)
}

func doUpdate(uri string, i interface{}) (err error) {
	url := uri + GetAccessToken()
	wxerr := new(WxErr)
	if err = util.PostJsonPtr(url, i, wxerr); err != nil {
		return
	}
	return wxerr.Error()
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

// GetGender 获取性别
func GetGender(s string) string {
	if s == "1" {
		return "男"
	}
	if s == "2" {
		return "女"
	}
	return "未定义"
}

var toUserReplacer = strings.NewReplacer("|", ",", "，", ",")

// GetToUser 获取acl所包含的所有用户
func GetToUser(acl interface{}) (touser string) {
	s1 := strings.TrimSpace(fmt.Sprint(acl))
	if strings.ToLower(s1) == "@all" {
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
	if strings.ToLower(acl) == "@all" {
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
