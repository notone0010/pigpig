// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// etcd
package etcd

import (
	"context"
	"strings"
	"sync"

	"github.com/notone/pigpig/pkg/log"
)

type serviceDiscover struct {
	ds *datastore
}

var (
	proxyMap map[string][]byte
	mu       sync.Mutex
)

func newDiscover(ds *datastore) *serviceDiscover {
	return &serviceDiscover{
		ds: ds,
	}
}

func (d *serviceDiscover) DiscoverStart(ctx context.Context, prefix string) error {
	return d.ds.Watch(ctx, prefix, d.SetProxyList, d.ModifyProxyList, d.DeleteProxyList)
}

func (d *serviceDiscover) RestartSession() error {
	return d.ds.RestartSession()
}

func (d *serviceDiscover) SetProxyList(ctx context.Context, key, addr []byte) {
	mu.Lock()
	defer mu.Unlock()
	if proxyMap == nil {
		proxyMap = make(map[string][]byte, 1)
	}
	proxyMap[string(key)] = addr
	log.Infof("proxy node list added %s", key)
}

func (d *serviceDiscover) ModifyProxyList(ctx context.Context, key, oldAddr, addr []byte) {
	mu.Lock()
	defer mu.Unlock()
	proxyMap[string(key)] = addr
	log.Infof("proxy node list modified %s", key)
}

func (d *serviceDiscover) DeleteProxyList(ctx context.Context, key []byte) {
	mu.Lock()
	defer mu.Unlock()
	delete(proxyMap, string(key))
	log.Infof("proxy node list deleted %s", key)
}

func (d *serviceDiscover) GetService() []string {
	mu.Lock()
	defer mu.Unlock()
	addrs := make([]string, 0)
	for k, v := range proxyMap {
		if !strings.Contains(k, "https") {
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (d *serviceDiscover) GetHttpsService() []string {
	mu.Lock()
	defer mu.Unlock()
	addrs := make([]string, 0)
	for k, v := range proxyMap {
		if strings.Contains(k, "https") {
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}


func (d *serviceDiscover) Close() error {
	return d.ds.Close()
}
