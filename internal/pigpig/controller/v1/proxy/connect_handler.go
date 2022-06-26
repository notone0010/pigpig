// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// proxy
package proxy

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

// ConnectHandler deal with http-request that method is 'connect'
func (p *ProxyController) ConnectHandler(w http.ResponseWriter, r *http.Request, targetHttps string) {
	// 1.将http method is connect 的请求发送到远端https服务器
	// 2.接收响应并转发给客户端
	targetAddress := net.JoinHostPort(p.LocalNetIFAddr, strconv.Itoa(p.SecureServingBindPort))
	if targetHttps != "" {
		targetAddress = targetHttps
	}
	destConn, err := net.DialTimeout("tcp", targetAddress, 60*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	_, _ = io.Copy(destination, source)
}
