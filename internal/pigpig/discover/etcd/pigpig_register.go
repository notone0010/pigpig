// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// etcd
package etcd

import (
	"context"
	"fmt"

	"github.com/notone/pigpig/pkg/log"
)

type register struct {
	ds *datastore
}

var (
	keyRegister = "register/%v/%v/%v"

	protoHttp  = "http"
	protoHttps = "https"

	globalClusterId string
	globalKey       string
	globalValue     string

	globalHttpsValue string
)

func newRegister(ds *datastore) *register {
	return &register{
		ds: ds,
	}
}

func (r *register) getKey(clusterId, proto, name string) string {
	return fmt.Sprintf(keyRegister, clusterId, proto, name)
}

// NewServiceRegister handle new service register into etcd for server discover
func (r *register) NewServiceRegister(ctx context.Context, clusterId, proto, key, value string) error {
	_key := r.getKey(clusterId, proto, key)
	if err := r.ds.PutSession(ctx, _key, value); err != nil {
		log.Errorf("machine ID %s failed put into etcd", key)
		return err
	}
	globalKey = key
	globalClusterId = clusterId
	switch proto {
	case protoHttp:
		globalValue = value
	case protoHttps:
		globalHttpsValue = value
	}
	log.Infof("machine ID %s put into etcd successful", key)
	return nil
}

func (r *register) RestartSession() error {
	return r.ds.RestartSession()
}

func (r *register) RecoveryServiceRegister(ctx context.Context) error {
	if globalKey != "" && globalValue != "" {
		httpKey := r.getKey(globalClusterId, "http", globalKey)
		if err := r.ds.PutSession(ctx, httpKey, globalValue); err != nil {
			log.Errorf("machine ID %s failed recover http service into etcd", globalKey)
			return err
		}
	}
	if globalKey != "" && globalHttpsValue != "" {
		httpsKey := r.getKey(globalClusterId, "https", globalKey)
		if err := r.ds.PutSession(ctx, httpsKey, globalHttpsValue); err != nil {
			log.Errorf("machine ID %s failed recover https service into etcd", globalKey)
			return err
		}
	}

	log.Infof("machine ID %s recover into etcd successful", globalKey)
	return nil
}

func (r *register) Close() error {
	ctx := context.TODO()
	if _, err := r.ds.cli.Revoke(ctx, r.ds.leaseID); err != nil {
		return err
	}
	return r.ds.Close()
}
