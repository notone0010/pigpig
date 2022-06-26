// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// proxy
package v1

import (
	"net/http"
	"time"

	"github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/internal/pigpig/transport"
	"github.com/notone/pigpig/pkg/log"
)

type ProxySrv interface {

	// PrepareRequest prepare request`s informations and options
	// PrepareRequest(r *http.Request) *dudu.ClientDetail

	// FetchRemoteResponse get response from remote web
	FetchRemoteResponse(c *dudu.Context)
}

type proxyService struct {
	transport transport.Factory

	// you can write your modify request and modify response function in here for the service handler call them
	// PluginChin []dudu.Plugin

	ModifyError dudu.ModifyError
}

// 1. 初始化客户端实例，监测请求头中的相关请求完成特定中间件的配置
// 2. 调度所有插件，汇总所有请求处理逻辑
// 3. 调用引擎逻辑，发起请求
// 4. 调度所有插件，汇总所有响应处理逻辑
// 5. 返回响应给上层接口

var _ ProxySrv = (*proxyService)(nil)

func newProxy(srv *service) *proxyService {
	return &proxyService{transport: srv.transport}
}

func (p *proxyService) FetchRemoteResponse(c *dudu.Context) {
	// create pool for prepare request object

	detail := c.RequestDetail

	sendTime := time.Now()
	resp, body, remoteErr := p.transport.GoRequest().FetchRemoteResponse(detail)
	if remoteErr != nil {
		log.Errorf("fetch remote response found error %s", remoteErr.Error())
		c.Errors = append(c.Errors, remoteErr)
		c.Abort()
		return
	}

	finalResponse := NewResponseDetail(resp, body)

	finalResponse.ElapsedTime = time.Now().Sub(sendTime)

	c.ResponseDetail = finalResponse
}

func NewResponseDetail(response *http.Response, body []byte) *dudu.ResponseDetail {
	responseDetail := &dudu.ResponseDetail{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Body:       body,
		RawBody:    response.Body,
		Response:   response,
	}
	return responseDetail
}
