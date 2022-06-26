// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// LB
package LB

import (
	"fmt"
	"sync"
	"testing"
)

func TestLB_RR(t *testing.T) {
	indexes := []int{0, 1, 3, 4}
	rr := RR{}
	resMap := make(map[int]int, 100)
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
			_index, _ := rr.SwitchTo(len(indexes))
			resMap[_index] = resMap[_index] + 1
			// fmt.Println(_index)
		}()
	}
	wg.Wait()
	for k, v := range resMap {
		fmt.Printf("key: %d -> value: %d\n", k, v)
	}
}
