package registry

type Config struct {
	Token       string // Token 令牌
	Namespace   string // 命名空间
	ServiceName string // 服务名称
	IP          string // 服务IP
	Port        int    // 服务端口
	Weight      int    // 权重
	Metadata    string // 元数据
	HeartBeat   int    // 健康上报间隔，单位：秒
}
