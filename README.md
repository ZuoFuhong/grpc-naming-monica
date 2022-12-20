## Monica naming SDK

gRPC [命名解析](https://github.com/grpc/grpc/blob/master/doc/naming.md)和[负载均衡](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md)的扩展实现，使用[monica](https://github.com/ZuoFuhong/monica) 注册中心.

### Usage

```shell
go get github.com/ZuoFuhong/grpc-naming-monica
```

服务提供方启动时，将服务地址注册到服务注册中心，同时定期报心跳到服务注册中心以表明服务的存活状态，相当于健康检查.

```go
regIns := NewRegistry(&Config{
    Token:       "18ee7064-3cdd-4ed5-a139-fd8d9add5847",
    Namespace:   "Test",
    ServiceName: "go_wallet_manage_svr",
    IP:          "127.0.0.1",
    Port:        1024,
    Weight:      100,
    Metadata:    "[]"
})
if err := regIns.Register(); err != nil {
    t.Fatal(err)
}
```

服务消费方要访问某个服务时，它通过 “服务发现组件” 向服务注册中心发出服务名称查询，名称将解析为一个或多个 IP 地址，同时缓存并定期刷新目标服务地址列表.

```go
conn, err := grpc.Dial("monica://Production/go_wallet_manage_svr")
if err != nil {
    t.Fatal(err)
}
```

服务消费方要访问某个服务时，根据配置的负载均衡策略选择一个目标服务地址，最后向目标服务发起请求。

```go
// gRPC 提供两种负载均衡策略 pick_first、round_robin, 默认的策略 pick_first
// 自定义实现 "加权轮询" 负载策略：weighted_round_robin
conn, err := grpc.Dial("monica://Test/go_wallet_manage_svr",
	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"weighted_round_robin"}`)
)
```

## License

This project is licensed under the [Apache 2.0 license](https://github.com/ZuoFuhong/grpc-naming-monica/blob/master/LICENSE).