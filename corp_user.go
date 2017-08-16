package wechat

import (
	"fmt"
	"log"
	"strings"

	"github.com/esap/wechat/util"
)

// WXAPI 企业号用户列表接口
const (
	WXAPI_GetUser     = WXAPI_ENT + "user/getuserinfo?access_token=%s&code=%s"
	WXAPI_GetUserInfo = WXAPI_ENT + "user/get?access_token=%s&userid=%s"
	WXAPI_UserList    = WXAPI_ENT + `user/list?access_token=%s&department_id=1&fetch_child=1&status=0`
	WXAPI_UserAdd     = WXAPI_ENT + `user/create?access_token=`
	WXAPI_UserUpdate  = WXAPI_ENT + `user/update?access_token=`
	WXAPI_UserDel     = WXAPI_ENT + `user/delete?access_token=`
)

// UserOauth 用户鉴权信息
type UserOauth struct {
	WxErr
	UserId   string
	DeviceId string
	OpenId   string
}

// GetUserOauth 通过code鉴权
func (s *Server) GetUserOauth(code string) (o UserOauth, err error) {
	url := fmt.Sprintf(WXAPI_GetUser, s.GetAccessToken(), code)
	if err = util.GetJson(url, &o); err != nil {
		return
	}
	if o.ErrCode != 0 {
		err = fmt.Errorf("GetUserId error : errcode=%v , errmsg=%v", o.ErrCode, o.ErrMsg)
	}
	return
}

// GetUserOauth 通过code鉴权
func GetUserOauth(code string) (userOauth UserOauth, err error) {
	return std.GetUserOauth(code)
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
	Telephone  string `json:"telephone,omitempty"`
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
func (s *Server) UserAdd(user *UserInfo) (err error) {
	return s.doUpdate(WXAPI_UserAdd, user)
}

// UserUpdate 添加用户
func (s *Server) UserUpdate(user *UserInfo) (err error) {
	return s.doUpdate(WXAPI_UserUpdate, user)
}

// UserDelete 删除用户
func (s *Server) UserDelete(user string) (err error) {
	e := new(WxErr)
	if err = util.GetJson(WXAPI_UserDel+s.GetAccessToken()+"&userid="+user, e); err != nil {
		return
	}
	return e.Error()
}

// GetUserInfo 通过userId获取用户信息
func (s *Server) GetUserInfo(userId string) (user UserInfo, err error) {
	url := fmt.Sprintf(WXAPI_GetUserInfo, s.GetAccessToken(), userId)
	if err = util.GetJson(url, &user); err != nil {
		return
	}
	if user.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", user.ErrCode, user.ErrMsg)
	}
	return
}

// GetUser 通过账号获取用户信息
func (s *Server) GetUser(userid string) *UserInfo {
	for _, v := range s.UserList.UserList {
		if v.UserId == userid {
			return &v
		}
	}
	return nil
}

// GetUserName 通过账号获取用户信息
func (s *Server) GetUserName(userid string) string {
	for _, v := range s.UserList.UserList {
		if v.UserId == userid {
			return v.Name
		}
	}
	return ""
}

// Users 用户列表
var Users userList

// UserList 用户列表
type userList struct {
	WxErr
	UserList []UserInfo
}

// SyncUserList 获取用户列表
func (s *Server) SyncUserList() (err error) {
	s.UserList, err = s.GetUserList()
	if err != nil {
		log.Println("同步通讯录失败:", err)
	}
	return
}

// GetUserList 获取用户列表
func (s *Server) GetUserList() (u userList, err error) {
	url := fmt.Sprintf(WXAPI_UserList, s.GetAccessToken())
	if err = util.GetJson(url, &u); err != nil {
		return
	}
	if u.ErrCode != 0 {
		err = fmt.Errorf("GetUserList error : errcode=%v , errmsg=%v", u.ErrCode, u.ErrMsg)
	}
	return
}

// GetUserNameList 获取用户列表
func (s *Server) GetUserNameList() (userlist []string) {
	userlist = make([]string, 0)
	for _, v := range s.UserList.UserList {
		userlist = append(userlist, v.UserId)
	}
	return
}

func (s *Server) doUpdate(uri string, i interface{}) (err error) {
	url := uri + s.GetAccessToken()
	wxerr := new(WxErr)
	if err = util.PostJsonPtr(url, i, wxerr); err != nil {
		return
	}
	return wxerr.Error()
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

var toUserReplacer = strings.NewReplacer(",", "|", "，", "|")

// GetToUser 获取acl所包含的所有用户
func (s *Server) GetToUser(acl interface{}) (touser string) {
	s1 := strings.TrimSpace(fmt.Sprint(acl))
	if strings.ToLower(s1) == "@all" {
		return "@all"
	}
	arr := strings.Split(toUserReplacer.Replace(s1), "|")
	for _, toUser := range arr {
		for _, v := range s.UserList.UserList {
			if s.CheckUserAcl(v.UserId, toUser) {
				touser += "|" + v.UserId
			}
		}
	}
	return strings.Trim(touser, "|")
}

// CheckUserAcl 测试权限，对比user的账号，姓名，手机，职位是否包含于acl
func (s *Server) CheckUserAcl(userid, acl string) bool {
	acl = strings.TrimSpace(acl)
	if acl == "" {
		return false
	}
	if strings.ToLower(acl) == "@all" {
		return true
	}
	acl = "," + strings.Replace(acl, "，", ",", -1) + ","
	u := s.GetUser(userid)
	if u == nil {
		return false
	}
	for _, dv := range u.Department {
		if strings.Contains(acl, ","+s.GetDeptName(dv)+",") {
			return true
		}
		if strings.Contains(acl, ","+s.GetDeptName(dv)+"/"+u.Position+",") {
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
