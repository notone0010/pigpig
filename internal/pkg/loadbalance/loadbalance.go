// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// loadbalance
package loadbalance

// LB type can choose among round-robin, hash and consistent-hash

var client Factory

type Factory interface {
	RR() RRLB
	Shuffle() ShuffleLB
	ConsistentHash() ConsistentHashLB
	// SwitchTo(policy LbType, alternativeIndexes []int) (int, error)
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
