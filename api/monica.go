package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	if _, err := unifiedPostRequest(registerUrl, bodyBytes); err != nil {
		return err
	}
	return nil
}

// Renew 服务刷新
func Renew(req *RenewReq) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if _, err := unifiedPostRequest(renewUrl, bodyBytes); err != nil {
		return err
	}
	return nil
}

// Deregister 注销实例
func Deregister(req *DeregisterReq) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if _, err := unifiedPostRequest(deregisterUrl, bodyBytes); err != nil {
		return err
	}
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
	// 目前不支持
	return []*InstanceNode{}, nil
}

// 统一请求
func unifiedPostRequest(reqUrl string, reqBody []byte) ([]byte, error) {
	resp, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed, statusCode = %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rspBody, nil
}
