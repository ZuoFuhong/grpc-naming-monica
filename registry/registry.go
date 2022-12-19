package registry

import (
	"fmt"
	"github.com/ZuoFuhong/grpc-naming-monica/api"
	"google.golang.org/grpc/grpclog"
	"time"
)

var logger = grpclog.Component("monica-registy")

const DefaultHeartBeat = 5 // 默认健康上报间隔

type Registry struct {
	cfg *Config
}

func NewRegistry(cfg *Config) *Registry {
	if cfg.HeartBeat == 0 {
		cfg.HeartBeat = DefaultHeartBeat
	}
	return &Registry{cfg: cfg}
}

// Register 服务注册
func (s *Registry) Register() error {
	err := api.Register(&api.RegisterReq{
		Token:       s.cfg.Token,
		Namespace:   s.cfg.Namespace,
		ServiceName: s.cfg.ServiceName,
		Node: &api.InstanceNode{
			IP:       s.cfg.IP,
			Port:     s.cfg.Port,
			Weight:   s.cfg.Weight,
			Metadata: s.cfg.Metadata,
		},
	})
	if err != nil {
		return fmt.Errorf("register error: %s", err.Error())
	}
	s.renew()
	return nil
}

// renew 服务更新
func (s *Registry) renew() {
	tick := time.Second * time.Duration(s.cfg.HeartBeat)
	go func() {
		for {
			err := api.Renew(&api.RenewReq{
				Token:       s.cfg.Token,
				Namespace:   s.cfg.Namespace,
				ServiceName: s.cfg.ServiceName,
				IP:          s.cfg.IP,
			})
			if err != nil {
				logger.Errorf("renew failed, err: %v", err)
			}
			time.Sleep(tick)
		}
	}()
}

// Deregister 服务注销
func (s *Registry) Deregister() error {
	err := api.Deregister(&api.DeregisterReq{
		Token:       s.cfg.Token,
		Namespace:   s.cfg.Namespace,
		ServiceName: s.cfg.ServiceName,
		IP:          s.cfg.IP,
	})
	return fmt.Errorf("deregister error: %s", err.Error())
}
