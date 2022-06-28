// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package loadbalance

// ConsistentHashLB consistent hash algorithm interface.
type ConsistentHashLB interface {
	Add(keys ...string)
	Get(key string) string
}
