// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// LB
package LB

import (
	"math/rand"
	"time"
)

// shuffle 算法当洗牌次数足够大时能够达到均等，否则算法会倾向于数组的0索引

func init() {
	rand.Seed(time.Now().UnixNano())
}

func shuffle(n int) []int {
	return rand.Perm(n)
}

type Shuffle struct {
	// you can set some options
}

func (s *Shuffle) SwitchTo(indexSize int) (int, error) {
	if indexSize <= 1 {
		return indexSize, nil
	}
	_indexes := shuffle(indexSize)
	return _indexes[0], nil
}
