package wechat

import (
	"github.com/esap/wechat/util"
)

const (
	// CorpAPICheckInGet 企业微信打开数据获取接口
	CorpAPICheckInGet = CorpAPI + "checkin/getcheckindata?access_token="
	// CorpCheckInAgentID  打卡AgentId
	CorpCheckInAgentID = 3010011
)

type (
	// dkDataReq 审批请求数据
	dkDataReq struct {
		OpenCheckInDataType int64    `json:"opencheckindatatype"`
		Starttime           int64    `json:"starttime"`
		Endtime             int64    `json:"endtime"`
		UseridList          []string `json:"useridlist"`
	}

	// DkDataRet 审批返回数据
	DkDataRet struct {
		WxErr  `json:"-"`
		Result []DkData `json:"checkindata"`
	}

	// DkData 审批数据
	DkData struct {
		Userid         string `json:"userid"`          // 用户id
		GroupName      string `json:"groupname"`       // 打卡规则名称
		CheckinType    string `json:"checkin_type"`    // 打卡类型
		ExceptionType  string `json:"exception_type"`  // 异常类型，如果有多个异常，以分号间隔
		CheckinTime    int64  `json:"checkin_time"`    // 打卡时间。UTC时间戳
		LocationTitle  string `json:"location_title"`  // 打卡地点title
		LocationDetail string `json:"location_detail"` // 打卡地点详情
		WifiName       string `json:"wifiname"`        // 打卡wifi名称
		Notes          string `json:"notes"`           // 打卡备注
		WifiMac        string `json:"wifimac"`         // 打卡的MAC地址/bssid
	}
)

// GetCheckIn 获取打卡数据,Namelist用户列表不超过100个。若用户超过100个，请分批获取
func (s *Server) GetCheckIn(opType, start, end int64, Namelist []string) (dkdata []DkData, err error) {
	url := CorpAPICheckInGet + s.GetAccessToken()
	data := new(DkDataRet)
	if err = util.PostJsonPtr(url, dkDataReq{opType, start, end, Namelist}, data); err != nil {
		return
	}
	if data.ErrCode != 0 {
		err = data.Error()
	}
	dkdata = data.Result
	return
}

// GetAllCheckIn 获取所有人的打卡数据
func (s *Server) GetAllCheckIn(opType, start, end int64) (dkdata []DkData, err error) {
	ul := s.GetUserIdList()
	l := len(ul)
	for i := 0; i < l; i += 100 {
		dk, e := s.GetCheckIn(opType, start, end, ul[i:util.Min(l, i+100)])
		if e != nil {
			err = e
			return
		}
		dkdata = append(dkdata, dk...)
	}
	return
}
