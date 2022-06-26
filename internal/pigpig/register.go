// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// pigpig
package pigpig

import (
	"github.com/notone/pigpig/internal/pigpig/controller/v1/proxy"
	"github.com/notone/pigpig/internal/pigpig/transport/xgorequest"
	"github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/internal/pkg/loadbalance/LB"
)

func registerProxyHandler(m *dudu.ProxyHttpMux) {
	installMiddleware(m)
	installController(m)
}

func installMiddleware(m *dudu.ProxyHttpMux) {

}

func installController(m *dudu.ProxyHttpMux) {
	// m.ProxyRequestHandler()
	// an transport will has initialized
	requestTransport := xgorequest.GetGorequestTransport()
	lb := LB.NewLB()
	proxyHandler := proxy.NewUserController(m, requestTransport, lb)

	proxyHandler.Plugins = m.GetHandlers()

	proxyHandler.InsecureServingBindPort = m.InsecureServingBindPort
	proxyHandler.SecureServingBindPort = m.SecureServingBindPort
	proxyHandler.LocalNetIFAddr = m.LocalNetIFAddr
	m.ProxyRequestHandler(proxyHandler.ServeHandle)
}
