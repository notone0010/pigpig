// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// LB
package LB

import "sync"

type RR struct {
}

var (
	index int
	mu    sync.Mutex
)

func (r *RR) SwitchTo(indexSize int) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	index++
	_index := index % indexSize
	return _index, nil
}
