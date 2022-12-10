package naming

import (
	"context"
	"fmt"
	"github.com/ZuoFuhong/grpc-naming-monica/api"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"log"
	"strings"
	"time"
)

func init() {
	resolver.Register(newMonicaBuilder())
}

var mlogger = grpclog.Component("monica-resolver")

const (
	// 同步实例列表的周期
	syncNSInterval = 5 * time.Second
)

type monicaBuilder struct{}

func newMonicaBuilder() resolver.Builder {
	return &monicaBuilder{}
}

func (m *monicaBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &monicaResolver{
		target: target,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
	}
	// 启动协程
	go r.watcher()
	return r, nil
}

func (m *monicaBuilder) Scheme() string {
	return "monica"
}

type monicaResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func (r *monicaResolver) ResolveNow(resolver.ResolveNowOptions) {
	log.Println("monica resolver resolve now")
}

func (r *monicaResolver) Close() {
	log.Println("monica resolver close")
	r.cancel()
}

// 轮询并更新服务的实例
func (r *monicaResolver) watcher() {
	r.updateState()
	ticker := time.NewTicker(syncNSInterval)
	for {
		select {
		case <-ticker.C:
			r.updateState()
		case <-r.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// 更新实例列表
func (r *monicaResolver) updateState() {
	instances := r.getInstances()
	newAddrs := make([]resolver.Address, 0)
	for _, ins := range instances {
		addr := resolver.Address{Addr: fmt.Sprintf("%s:%d", ins.IP, ins.Port)}
		// 通过属性存储权重
		addr = SetAddrInfo(addr, AddrInfo{
			Weight: ins.Weight,
		})
		newAddrs = append(newAddrs, addr)
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: newAddrs})
}

// 获取服务可用的实例
func (r *monicaResolver) getInstances() []*api.InstanceNode {
	ns := r.target.URL.Host
	sname := strings.TrimPrefix(r.target.URL.Path, "/")
	// 调用 Monica API
	nodes, err := api.Fetch(ns, sname)
	if err != nil {
		mlogger.Errorf("fetch instance error: %v", err)
		return []*api.InstanceNode{}
	}
	return nodes
}
