package naming

import (
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"sync"
)

// Name is the name of weighted_round_robin balancer.
const Name = "weighted_round_robin"

var logger = grpclog.Component("weightedroundrobin")

func init() {
	balancer.Register(newBuilder())
}

// newBuilder creates a new weightedroundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &wrrPickerBuiler{}, base.Config{HealthCheck: true})
}

type attributeKey struct{}

type AddrInfo struct {
	Weight int
}

func SetAddrInfo(addr resolver.Address, addrInfo AddrInfo) resolver.Address {
	addr.Attributes = attributes.New(attributeKey{}, addrInfo)
	return addr
}

func GetAddrInfo(addr resolver.Address) AddrInfo {
	return addr.Attributes.Value(attributeKey{}).(AddrInfo)
}

type wrrPickerBuiler struct{}

func (*wrrPickerBuiler) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("roundrobinPicker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	nodes := make([]*Node, 0, len(info.ReadySCs))
	for subConn, addr := range info.ReadySCs {
		addrInfo := GetAddrInfo(addr.Address)
		node := &Node{
			weight:          addrInfo.Weight,
			currentWeight:   addrInfo.Weight,
			effectiveWeight: addrInfo.Weight,
			addr:            addr.Address.Addr,
			conn:            subConn,
		}
		nodes = append(nodes, node)
	}
	return &wrrPicker{
		nodes: nodes,
	}
}

type Node struct {
	weight          int
	currentWeight   int
	effectiveWeight int
	addr            string
	conn            balancer.SubConn
}

type wrrPicker struct {
	nodes    []*Node
	mu       sync.Mutex
	curIndex int
}

// Pick 加权轮询
func (p *wrrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	totalWeight := 0
	var maxWeightNode *Node
	for index, node := range p.nodes {
		totalWeight += node.effectiveWeight
		node.currentWeight += node.effectiveWeight
		if maxWeightNode == nil || maxWeightNode.currentWeight < node.currentWeight {
			maxWeightNode = node
			p.curIndex = index
		}
	}
	maxWeightNode.currentWeight -= totalWeight
	p.mu.Unlock()
	return balancer.PickResult{SubConn: maxWeightNode.conn}, nil
}
