// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// proxy
package proxy

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/marmotedu/errors"
	"github.com/notone/pigpig/internal/pigpig/discover"
	srvv1 "github.com/notone/pigpig/internal/pigpig/service/v1"
	"github.com/notone/pigpig/internal/pigpig/transport"
	"github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/internal/pkg/loadbalance"
	"github.com/notone/pigpig/internal/pkg/loadbalance/LB"
	"github.com/notone/pigpig/pkg/core"
	"github.com/notone/pigpig/pkg/log"
)

type ProxyController struct {
	srv srvv1.Service

	engine *dudu.ProxyHttpMux

	Lb LB.LB

	InsecureServingBindPort int
	SecureServingBindPort   int

	// LocalNetIFAddr is the network interface address the current local machine
	LocalNetIFAddr string


	Plugins dudu.HandlersChain

	handleErrorFunc dudu.HandlerErrorFunc

	// connPool connect pool
	connPool sync.Pool

	// userPool
	userPool sync.Pool
}

// NewUserController creates a user handler.
func NewUserController(engine *dudu.ProxyHttpMux,
	transport transport.Factory,
	lb LB.LB,
) *ProxyController {
	controller := &ProxyController{
		engine: engine,
		srv:    srvv1.NewService(transport),
		Lb:     lb,
	}
	controller.userPool.New = func() interface{} {
		return controller.allocateContext()
	}
	return controller
}

// ChooseHandler select handler by http method
func (p *ProxyController) ServeHandle(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(dudu.InternalHeaderFullPath) == "" {
		r.Header.Add(dudu.InternalHeaderFullPath, r.RequestURI)
	}

	// handle reverse proxy
	if p.engine.Cluster.Enable &&
		strings.Compare(p.engine.Cluster.Role, dudu.RoleMaster) == 0 &&
		r.TLS == nil {
		// load balance
		var (
			nextIndex   int
			serviceList []string
		)

		// will use etcd to handle config`s update
		if r.Method == http.MethodConnect {
			serviceList = discover.Client().Discover().GetHttpsService()
		} else {
			serviceList = discover.Client().Discover().GetService()
		}
		log.Infof("online service list ---> %s", serviceList)

		if p.engine.Cluster.LoadPolicy == loadbalance.RR {
			nextIndex, _ = p.Lb.RR().SwitchTo(len(serviceList))
		} else if p.engine.Cluster.LoadPolicy == loadbalance.Shuffle {
			nextIndex, _ = p.Lb.Shuffle().SwitchTo(len(serviceList))
		} else {
			// default choose round-robin
			nextIndex, _ = p.Lb.RR().SwitchTo(len(serviceList))
		}
		targetHost := serviceList[nextIndex]

		if r.Method == http.MethodConnect {
			if targetHost == net.JoinHostPort(p.LocalNetIFAddr, strconv.Itoa(p.SecureServingBindPort)) {
				p.GoHandle(w, r)
				return
			}
			log.Debugf("Moved to %s method: CONNECT host: %s ", r.Host, targetHost)
			p.ConnectHandler(w, r, targetHost)
			return
		} else {
			if targetHost == net.JoinHostPort(p.LocalNetIFAddr, strconv.Itoa(p.InsecureServingBindPort)) {
				p.GoHandle(w, r)
				return
			}

			targetURL := dudu.HttpPrefix + targetHost

			reverseFunc := GoReverseProxy(targetURL)
			if reverseFunc == nil {
				reverseErr := errors.New("failed to get reverse proxy function")
				core.WriteResponse(w, r, reverseErr, nil)
				return
			}
			log.Debugf("Moved to %s method: %s host: %s ", r.Method, r.Host, targetHost)

			reverseFunc(w, r)
			return
		}

	}
	p.GoHandle(w, r)
}

func (p *ProxyController) allocateContext() *dudu.Context {
	return dudu.NewContext(p.engine)
}

func (p *ProxyController) GoHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		log.Infof("client: %s | received https CONNECT request %s", r.RemoteAddr, r.Host)
		p.ConnectHandler(w, r, "")
		return
	}
	c := p.userPool.Get().(*dudu.Context)
	c.GetContextObj(w, r)
	log.Infof("client: %s received %s request to %s %s", r.RemoteAddr, r.Method, r.Host, r.URL.Path)
	p.UserRequestHandler(c)

	p.userPool.Put(c)
}

func (p *ProxyController) Use(plugins ...dudu.HandlerFunc) {
	p.Plugins = append(p.Plugins, plugins...)
}
