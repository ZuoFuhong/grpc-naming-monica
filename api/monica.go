package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net/http"
)

var logger = grpclog.Component("monica-api")

const (
	registerUrl   = "http://monica.oa.com:1024/api/v1/register"
	renewUrl      = "http://monica.oa.com:1024/api/v1/renew"
	deregisterUrl = "http://monica.oa.com:1024/api/v1/deregister"
	fetchUrl      = "http://monica.oa.com:1024/api/v1/fetch?ns=%s&sname=%s"
	pollUrl       = "http://monica.oa.com:1024/api/v1/poll?ns=%s&sname=%s"
)

// Register 服务注册
func Register(req *RegisterReq) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(registerUrl, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Infof("monica.register resp: %v", string(respBytes))
	return nil
}

// Renew 服务刷新
func Renew(req *RenewReq) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(renewUrl, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Infof("monica.renew resp: %v", string(respBytes))
	return nil
}

// Deregister 注销实例
func Deregister(req *DeregisterReq) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(deregisterUrl, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Infof("monica.deregister resp: %v", string(respBytes))
	return nil
}

// Fetch 获取服务实例
func Fetch(ns, sname string) ([]*InstanceNode, error) {
	resp, err := http.Get(fmt.Sprintf(fetchUrl, ns, sname))
	if err != nil {
		return nil, fmt.Errorf("monica.fetch error: %s", err.Error())
	}
	fetchRsp := new(FetchResp)
	if err := json.NewDecoder(resp.Body).Decode(fetchRsp); err != nil {
		return nil, fmt.Errorf("monica.fetch error: %s", err.Error())
	}
	return fetchRsp.Data, nil
}

// Poll 获取服务实例（长轮询）
func Poll(ns, sname string) ([]*InstanceNode, error) {
	resp, err := http.Get(fmt.Sprintf(pollUrl, ns, sname))
	if err != nil {
		return nil, err
	}
	fetchRsp := new(FetchResp)
	if err := json.NewDecoder(resp.Body).Decode(fetchRsp); err != nil {
		return nil, err
	}
	return fetchRsp.Data, nil
}
