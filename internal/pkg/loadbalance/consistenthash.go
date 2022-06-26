// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// loadbalance
package loadbalance

type ConsistentHashLB interface {
	Add(keys ...string)
	Get(key string) string
}
