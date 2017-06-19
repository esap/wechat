package corp

import (
	"log"

	"github.com/esap/wechat"
	"github.com/esap/wechat/util"
)

const (
	// WXAPI_GetApproval 企业号审批数据获取接口
	WXAPI_GetApproval = wechat.WXAPI_ENT + "corp/getapprovaldata?access_token="
	// Corp_Approval_agentId 审批AgentId
	Corp_Approval_agentId = 3010040
)

type (
	// spDataReq 审批请求数据
	spDataReq struct {
		Starttime int64  `json:"starttime"`
		Endtime   int64  `json:"endtime"`
		NextSpNum string `json:"next_spnum,omitempty"`
	}

	// SpDataRet 审批返回数据
	SpDataRet struct {
		wechat.WxErr `json:"-"`
		Count        int64    `json:"count"`
		Total        int64    `json:"total"`
		NextSpnum    int64    `json:"next_spnum"`
		Data         []SpData `json:"data""`
	}

	// SpData 审批数据
	SpData struct {
		Spname       string    `json:"spname"`        // 审批名称(请假，报销，自定义审批名称)
		ApplyName    string    `json:"apply_name"`    // 申请人姓名
		ApplyOrg     string    `json:"apply_org"`     // 申请人部门
		ApprovalName []string  `json:"approval_name"` // 审批人姓名
		NotifyName   []string  `json:"notify_name"`   // 抄送人姓名
		SpStatus     int64     `json:"sp_status"`     // 审批状态：1审批中；2 已通过；3已驳回；4已取消
		SpNum        int64     `json:"sp_num"`        // 审批单号
		Leave        Leave     `json:"leave"`         // 请假类型
		Expense      Expense   `json:"expense"`       // 报销类型
		Comm         ApplyData `json:"comm"`          // 自定义类型
	}

	// Leave 请假
	Leave struct {
		Timeunit  int64  `json:"timeunit"`   // 请假时间单位：0半天；1小时
		LeaveType int64  `json:"leave_type"` // 请假类型：1年假；2事假；3病假；4调休假；5婚假；6产假；7陪产假；8其他
		StartTime int64  `json:"start_time"` // 请假开始时间，unix时间
		EndTime   int64  `json:"end_time"`   // 请假结束时间，unix时间
		Duration  int64  `json:"duration"`   // 请假时长，单位小时
		Reason    string `json:"reason"`     // 请假事由
	}

	// Expense 报销
	Expense struct {
		ExpenseType int64         `json:"expense_type"` // 报销类型：1差旅费；2交通费；3招待费；4其他报销
		Reason      string        `json:"reason"`       // 报销事由
		Item        []ExpenseItem `json:"item"`         // 报销明细
	}

	// ExpenseItem 报销明细
	ExpenseItem struct {
		ExpenseitemType int64  `json:"expenseitem_type"` // 费用类型：1飞机票；2火车票；3的士费；4住宿费；5餐饮费；6礼品费；7活动费；8通讯费；9补助；10其他
		Time            int64  `json:"time"`             // 发生时间，unix时间
		Sums            int64  `json:"sums"`             // 费用金额，单位元
		Reason          string `json:"reason"`           // 明细事由
	}

	// ApplyData 自定义数据
	ApplyData struct {
		Data string `json:"apply_data"` // 自定义审批申请的单据数据
	}

	// ApplyField 自定义字段
	ApplyField struct {
		Title string      `json:"title"` // 类目名
		Type  string      `json:"type"`  // 类目类型【 text: "文本", textarea: "多行文本", number: "数字", date: "日期", datehour: "日期+时间",  select: "选择框" 】
		Value interface{} `json:"value"` // 填写的内容，只有Type是图片时，value是一个数组，数据示例如下方所示；
	}
)

// GetApproval 获取审批数据
func GetApproval(start, end int64, nextNum string) (spData *SpDataRet, err error) {
	at, err := wechat.GetAgentAccessToken(Corp_Approval_agentId)
	if err != nil {
		return nil, err
	}
	url := WXAPI_GetApproval + at
	spData = new(SpDataRet)
	if err = util.PostJsonPtr(url, spDataReq{start, end, nextNum}, spData); err != nil {
		log.Println("PostJsonPtr err:", err)
		return
	}
	if spData.ErrCode != 0 {
		err = spData.Error()
	}
	return
}
