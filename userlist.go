package wechat

import (
	"fmt"

	"github.com/esap/wechat/util"
)

// WXAPI_USERLIST 企业号用户列表接口
const (
	WXAPI_USERLIST = WXAPI_ENT + `user/list?access_token=%s&department_id=1&fetch_child=1&status=0`
)

type UserList struct {
	wxErr
	UserList []struct {
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
}

func GetUserList() (userList UserList, err error) {
	url := fmt.Sprintf(WXAPI_USERLIST, GetAccessToken())
	if err = util.HttpGetJson(url, &userList); err != nil {
		return
	}
	if userList.ErrCode > 0 {
		err = fmt.Errorf("MediaUpload error : errcode=%v , errmsg=%v", userList.ErrCode, userList.ErrMsg)
	}
	return
}
