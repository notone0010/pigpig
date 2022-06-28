// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lb

import "sync"

// RR round-robin.
type RR struct{}

var (
	index int
	mu    sync.Mutex
)

// SwitchTo choose next.
func (r *RR) SwitchTo(indexSize int) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	index++
	_index := index % indexSize

	return _index, nil
}
