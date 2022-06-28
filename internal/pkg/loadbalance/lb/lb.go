// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// LB
package LB

type LBType string

type LB struct {
}

func NewLB() LB {
	return LB{}
}

func (l *LB) Shuffle() *Shuffle {
	return &Shuffle{}
}

func (l *LB) RR() *RR {
	return &RR{}
}
