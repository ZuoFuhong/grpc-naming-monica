package api

// RegisterReq 服务注册
type RegisterReq struct {
	Token       string        `json:"token"`        // Token 令牌
	Namespace   string        `json:"namespace"`    // 命名空间
	ServiceName string        `json:"service_name"` // 服务名称
	Node        *InstanceNode `json:"node"`         // 实例节点
}

// RenewReq 服务更新
type RenewReq struct {
	Token       string `json:"token"`        // Token 令牌
	Namespace   string `json:"namespace"`    // 命名空间
	ServiceName string `json:"service_name"` // 服务名称
	IP          string `json:"ip"`           // 实例IP
}

// DeregisterReq 服务注销
type DeregisterReq struct {
	Token       string `json:"token"`        // Token 令牌
	Namespace   string `json:"namespace"`    // 命名空间
	ServiceName string `json:"service_name"` // 服务名称
	IP          string `json:"ip"`           // 实例IP
}
