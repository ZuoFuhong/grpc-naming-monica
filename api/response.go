package api

type InstanceNode struct {
	IP       string `json:"ip"`       // 实例IP
	Port     int    `json:"port"`     // 端口
	Weight   int    `json:"weight"`   // 权重
	Metadata string `json:"metadata"` // 元数据
}

type FetchResp struct {
	Retcode int             `json:"retcode"`
	Errmsg  string          `json:"errmsg"`
	Data    []*InstanceNode `json:"data"`
}
