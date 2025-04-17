package mymodals

// 定义上行数据结构登陆请求
type AccountServiceRequest struct {
	Type string                 `json:"type"`
	Tag  string                 `json:"tag"`
	Body map[string]interface{} `json:"Body"`
}

type AccountLoginBody struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	IP          string `json:"ip"`
	Captcha     string `json:"captcha"`
	Proving     string `json:"proving"`
	MachineCode string `json:"machineCode"`
	ClientType  string `json:"client_type"`
}

// 定义下行数据结构
type AccountServiceResponse struct {
	Type   string `json:"type"`
	Tag    string `json:"tag"`
	Result string `json:"result"`
	Body   struct {
		Length int    `json:"length"`
		Detail string `json:"detail,omitempty"` // 可选字段，当result为error时使用
	} `json:"body"`
}

// 定义下行数据结构获取用户账号信息
type GetAccountResponse struct {
	Type string      `json:"type"`
	Tag  string      `json:"tag"`
	Body AccountBody `json:"body"`
}

type AccountBody struct {
	AgentUserID      int    `json:"agent_user_id"`
	AvatarNo         int    `json:"avatar_no"`
	AwardValue       int    `json:"award_value"`
	BagMoney         int    `json:"bag_money"`
	BankMoney        int    `json:"bank_money"`
	BankValidate     int    `json:"bank_validate"`
	Birthday         string `json:"birthday"`
	CertNo           string `json:"cert_no"`
	CertType         int    `json:"cert_type"`
	ChickTask        int    `json:"chick_task"`
	DollarMoney      int    `json:"dollar_money"`
	Email            string `json:"email"`
	ExpValue         int    `json:"exp_value"`
	GiftCheck        bool   `json:"gift_check"`
	Happycard        int    `json:"happycard"`
	InviterID        string `json:"inviter_id"`
	Level            int    `json:"level"`
	Locked           int    `json:"locked"`
	Nickname         string `json:"nickname"`
	OnlineTime       int    `json:"online_time"`
	Password         string `json:"password"` // 注意：在实际应用中不应返回密码
	PayMoney         int    `json:"pay_money"`
	PromoDollar      int    `json:"promo_dollar"`
	QQ               string `json:"qq"`
	Realname         string `json:"realname"`
	Sex              int    `json:"sex"`
	Telephone        string `json:"telephone"`
	UserID           int    `json:"userid"`
	Username         string `json:"username"`
	ValidateCert     int    `json:"validate_cert"`
	ValidateEmail    int    `json:"validate_email"`
	ValidateMobile   int    `json:"validate_mobile"`
	ValidateQQ       int    `json:"validate_qq"`
	ValidateRealname int    `json:"validate_realname"`
	ValidateVideo    int    `json:"validate_video"`
	VipLevel         int    `json:"vip_level"`
	WinDollarSum     int    `json:"win_dollar_sum"`
	GoldMoney        int    `json:"gold_money"`
	NextExp          int    `json:"next_exp"`
	BankCard         string `json:"bank_card"`
	Gfname           string `json:"gfname"`
	Wfname           string `json:"wfname"`
	Gfname1          string `json:"gfname1"`
	Wfname1          string `json:"wfname1"`
}

type reqRequest struct {
	Type string                 `json:"type"`
	Tag  string                 `json:"tag"`
	Body map[string]interface{} `json:"body"` // 动态 JSON 对象
}
