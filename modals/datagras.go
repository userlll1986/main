package mymodals

// 定义上行数据结构
type AccountServiceRequest struct {
	Type        string `json:"type"`
	Tag         string `json:"tag"`
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
