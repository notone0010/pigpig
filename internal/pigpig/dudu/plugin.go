// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dudu

// Plugin 代理平台处理请求插件.
type Plugin interface {
	// BeforeSendRequest will handle a request before send the request
	BeforeSendRequest(c *Context)

	// BeforeSendResponse will handle a request when received a response from remote
	BeforeSendResponse(c *Context)

	// OnConnectError will handle a request when the request connect remote with a error
	// OnConnectError()
}

// ModifyError modify error response when the server encounter error.
type ModifyError func(request *RequestDetail, err error) *ResponseDetail

//
