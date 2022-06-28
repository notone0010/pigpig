// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lb

// LB load balance object.
type LB struct{}

// NewLB return new load balance object.
func NewLB() LB {
	return LB{}
}

// Shuffle returns shuffle algorithm.
func (l *LB) Shuffle() *Shuffle {
	return &Shuffle{}
}

// RR returns rr algorithm.
func (l *LB) RR() *RR {
	return &RR{}
}
