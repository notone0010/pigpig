// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notone0010/pigpig/internal/pigpig/discover"
	"github.com/notone0010/pigpig/internal/pigpig/dudu"
	"github.com/notone0010/pigpig/internal/pkg/plugins"
	"github.com/notone0010/pigpig/pkg/core"
	"github.com/notone0010/pigpig/pkg/storage"
	"github.com/notone0010/pigpig/pkg/util/infoutil"
	"golang.org/x/sync/errgroup"

	"github.com/notone0010/pigpig/internal/pkg/middleware"
	"github.com/notone0010/pigpig/pkg/log"

	"github.com/marmotedu/component-base/pkg/version"
)

// RedisKeyPrefix redis keys uniform prefix.
const RedisKeyPrefix = "pigpig-service-discover-"

// GenericProxyServer contains state for an iam api server.
// type GenericProxyServer gin.Engine.
type GenericProxyServer struct {
	middlewares []string
	// SecureServingInfo holds configuration of the TLS server.
	SecureServingInfo *SecureServingInfo

	// InsecureServingInfo holds configuration of the insecure HTTP server.
	InsecureServingInfo *InsecureServingInfo

	plugins []string

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	healthz         bool
	enableMetrics   bool
	enableProfiling bool

	Cluster dudu.Cluster

	DistributedNode *infoutil.DistributedId
	// wrapper for gin.Engine

	MachineId int64

	LocalNetIFAddr string

	httpsAddress string
	httpAddress  string

	Engine                       *dudu.ProxyHttpMux
	insecureServer, secureServer *http.Server
}

func initGenericProxyServer(s *GenericProxyServer) {
	// do some setup
	// s.GET(path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
	s.InstallPlugins()
}

// InstallAPIs install generic apis.
func (s *GenericProxyServer) InstallAPIs() {
	// install healthz handler
	if s.healthz {
		s.Engine.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			core.WriteResponse(w, r, nil, map[string]string{"status": "ok"})
		})
	}

	// // install metric handler
	// if s.enableMetrics {
	// 	prometheus := ginprometheus.NewPrometheus("gin")
	// 	prometheus.Use(s.Engine)
	// }

	// install pprof handler
	// if s.enableProfiling {
	// 	pprof.Register(s.Engine)
	// }

	s.Engine.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		core.WriteResponse(w, r, nil, version.Get())
	})

	// s.Engine.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	core.WriteResponse(w, r, nil, map[string]string{"message": "Hello, World!"})
	// })
}

// Setup do some setup work for gin transport.
func (s *GenericProxyServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares install generic middlewares.
func (s *GenericProxyServer) InstallMiddlewares() {
	// necessary middlewares
	// 自定义代理平台处理中间件
	// s.Engine.Use(middleware.Recovery())
	// s.Engine.Use(middleware.Logger())
	// s.Engine.Use(middleware.WithLogger)

	// install custom middlewares
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)

			continue
		}

		log.Infof("install middleware: %s", m)
		log.Infof("%s", mw)
		s.Engine.Use(mw)
	}
}

/*
// preparedGenericProxyServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedGenericProxyServer struct {
	*GenericProxyServer
}.
*/

// InstallPlugins install plugin.md.
func (s *GenericProxyServer) InstallPlugins() {
	po := plugins.NewPluginsOptions(s.plugins)
	po.LoadPlugins()
	if len(po.Plugins) > 0 {
		s.Engine.Use(po.Plugins...)
	}
}

// Run spawns the http server. It only returns when the port cannot be listened on initially.
func (s *GenericProxyServer) Run() error {
	// For scalability, use custom HTTP configuration mode here
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s.Engine,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,

	}

	// For scalability, use custom HTTP configuration mode here
	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s.Engine,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
		TLSConfig: GetTlsConfig(s.SecureServingInfo.CertKey),
	}

	// GetCertificate 处理证书，如果不存在则导入默认证书并修改commonName 完成调用

	var eg errgroup.Group

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}

		log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())

		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	// it will general machine id
	if err := s.GetMachineId(); err != nil {
		log.Fatal(err.Error())
	}

	if s.Cluster.Enable {
		s.setupCluster()
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (s *GenericProxyServer) setupCluster() {
	ctx := context.Background()

	distributedNode, err := infoutil.NewDistributedId(s.MachineId)
	if err != nil {
		log.Fatal(err.Error())
	}
	s.DistributedNode = distributedNode

	if s.Cluster.Role == "master" {
		s.Cluster.ClusterId = s.DistributedNode.Generate().String()

		if err := discover.Client().Info().Put(ctx, s.Cluster.Name, s.Cluster.ClusterId); err != nil {
			log.Fatalf("the server failed to put cluster master`s information into etcd err: %s", err.Error())
		}
		log.Infof("the server successfully put cluster master`s information into etcd cluster id: %s", s.Cluster.ClusterId)

		if err := discover.Client().Discover().DiscoverStart(ctx, "register/"+s.Cluster.ClusterId); err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("the server successfully watch %s", s.Cluster.Name)
		// master will sleep 1 second so that the service-register can put into etcd successful
		time.Sleep(2 * time.Second)
	}

	if s.Cluster.ClusterId == "" {
		clusterId, err := discover.Client().Info().Get(ctx, s.Cluster.Name)
		if err != nil {
			log.Fatalf("failed to join %s cluster error: %s", s.Cluster.Name, err.Error())
		}
		s.Cluster.ClusterId = string(clusterId)
	}

	if s.Cluster.IsMasterHandle || s.Cluster.Role == "slave" {
		// initial current node information into etcd

		s.httpAddress = net.JoinHostPort(s.LocalNetIFAddr, strconv.Itoa(s.InsecureServingInfo.BindPort))

		s.httpsAddress = net.JoinHostPort(s.LocalNetIFAddr, strconv.Itoa(s.SecureServingInfo.BindPort))

		if err := discover.Client().Register().NewServiceRegister(ctx, s.Cluster.ClusterId, "http", strconv.FormatInt(s.MachineId, 10), s.httpAddress); err != nil {
			log.Fatal(err.Error())
		}
		if err := discover.Client().Register().NewServiceRegister(ctx, s.Cluster.ClusterId, "https", strconv.FormatInt(s.MachineId, 10), s.httpsAddress); err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("the current register into etcd successful")
	}

	log.Infof("server initial the cluster is success")
}

// Close graceful shutdown the api server.
func (s *GenericProxyServer) Close() {
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}

	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}
}

// ping pings the http server to make sure the router is working.
func (s *GenericProxyServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		// Change NewRequest to NewRequestWithContext and pass context it
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// Ping the server by sending a GET request to `/healthz`.

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}

// GetMachineId get current server machine id.
func (s *GenericProxyServer) GetMachineId() error {
	redisStorage := storage.RedisCluster{KeyPrefix: RedisKeyPrefix}

	machineIdKey := RedisKeyPrefix + "increase-machine-id"
	// time.Sleep(5*time.Second)
	for i := 0; i < 5; i++ {
		value := redisStorage.IncreaseWithExpire(machineIdKey, -1)
		if value != 0 {
			s.MachineId = value
			log.Infof("redis key: increase-machine-id ---> %v", value)

			return nil
		}
		log.Debugf("be failed to connect redis is trying to retry")
		time.Sleep(1 * time.Second)
	}

	return errors.New("failed to get machine id")
}
