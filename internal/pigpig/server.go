// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pigpig

import (
	"context"
	"fmt"
	"time"

	"github.com/marmotedu/errors"
	"github.com/notone0010/pigpig/internal/pigpig/config"
	"github.com/notone0010/pigpig/internal/pigpig/discover"
	"github.com/notone0010/pigpig/internal/pigpig/discover/etcd"

	// "github.com/notone0010/pigpig/internal/pigpig/transport".
	genericoptions "github.com/notone0010/pigpig/internal/pkg/options"
	genericproxyserver "github.com/notone0010/pigpig/internal/pkg/server"
	"github.com/notone0010/pigpig/pkg/log"
	"github.com/notone0010/pigpig/pkg/shutdown"
	"github.com/notone0010/pigpig/pkg/shutdown/shutdownmanagers/posixsignal"
	"github.com/notone0010/pigpig/pkg/storage"
)

type proxyServer struct {
	gs                 *shutdown.GracefulShutdown
	redisOptions       *genericoptions.RedisOptions
	genericProxyServer *genericproxyserver.GenericProxyServer
}

type preparedProxyServer struct {
	*proxyServer
}

// ExtraConfig defines extra configuration for the iam-proxyServer.
type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   genericoptions.GeneratableKeyCert
	mysqlOptions *genericoptions.MySQLOptions
}

func createProxyServer(cfg *config.Config) (*proxyServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	if cfg.GenericServerRunOptions.Cluster.Enable && cfg.GenericServerRunOptions.Cluster.Role == "" {
		err := fmt.Errorf("config`s 'server.cluster' must contain role as 'server.cluster.enable' is ture")

		return nil, errors.New(err.Error())
	}
	if cfg.GenericServerRunOptions.Cluster.Enable {
		discoverIns, err := etcd.GetEtcdFactoryOr(cfg.EtcdOptions, func() {
			instance, err := etcd.GetEtcdFactoryOr(nil, nil)
			if err != nil {
				log.Errorf("etcd on keepalive failure function found an error %s", err.Error())

				return
			}
			err = instance.Register().RestartSession()
			if err != nil {
				log.Errorf("etcd failed to restart session error ---> %s", err.Error())
			}
			log.Infof("etcd restart session successful")

			for i := 1; i <= 5; i++ {
				err := instance.Register().RecoveryServiceRegister(context.TODO())
				if err == nil {
					return
				}
				log.Errorf(err.Error())
				time.Sleep(1 * time.Second)
			}
		})
		if err != nil {
			return nil, err
		}
		discover.SetClient(discoverIns)
	}
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &proxyServer{
		gs:                 gs,
		redisOptions:       cfg.RedisOptions,
		genericProxyServer: genericServer,
	}

	return server, nil
}

func (s *proxyServer) PrepareRun() preparedProxyServer {
	registerProxyHandler(s.genericProxyServer.Engine)

	s.initRedisStore()

	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		// graceful shutdown any processing server
		s.genericProxyServer.Close()
		etcdIns, _ := etcd.GetEtcdFactoryOr(nil, nil)
		if etcdIns != nil {
			_ = etcdIns.Close()
		}

		return nil
	}))

	// 完成requestHandler的注册和初始化

	return preparedProxyServer{s}
}

func (s preparedProxyServer) Run() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericProxyServer.Run()
}

type completedExtraConfig struct {
	*ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericproxyserver.Config, lastErr error) {
	genericConfig = genericproxyserver.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func (s *proxyServer) initRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		cancel()

		return nil
	}))

	_config := &storage.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}

	// try to connect to redis
	go storage.ConnectToRedis(ctx, _config)
}
