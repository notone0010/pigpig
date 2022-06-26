// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// etcd
package etcd

import "context"

type info struct {
	ds *datastore
}

func newInfo(ds *datastore) *info {
	return &info{
		ds: ds,
	}
}

func (i *info) Get(ctx context.Context, key string) ([]byte, error) {
	return i.ds.Get(ctx, key)
}

func (i *info) Put(ctx context.Context, key, value string) error {
	return i.ds.Put(ctx, key, value)
}

func (i *info) Close() error {
	return i.ds.Close()
}
