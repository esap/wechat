package corp

import (
	"github.com/esap/wechat"
	"github.com/esap/wechat/util"
)

const (
	// WXAPI_CheckIn 企业号打开数据获取接口
	WXAPI_CheckIn = wechat.WXAPI_ENT + "checkin/getcheckindata?access_token="
	// Corp_GetCheckIn_agentId 打卡AgentId
	Corp_GetCheckIn_agentId = 3010011
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
		wechat.WxErr `json:"-"`
		Result       []DkData `json:"checkindata""`
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

// GetCheckIn 获取打卡数据
func GetCheckIn(opType, start, end int64) (dkdata []DkData, err error) {
	at, err := wechat.GetAgentAccessToken(Corp_GetCheckIn_agentId)
	if err != nil {
		return nil, err
	}
	url := WXAPI_CheckIn + at
	data := new(DkDataRet)
	if err = util.PostJsonPtr(url, dkDataReq{opType, start, end, wechat.GetUserNameList()}, data); err != nil {
		return
	}
	if data.ErrCode != 0 {
		err = data.Error()
	}
	dkdata = data.Result
	return
}
