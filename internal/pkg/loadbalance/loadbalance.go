// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package loadbalance

// lb type can choose among round-robin, hash and consistent-hash

var client Factory

// Factory load balance factory.
type Factory interface {
	RR() RRLB
	Shuffle() ShuffleLB
	ConsistentHash() ConsistentHashLB
	// SwitchTo(policy LbType, alternativeIndexes []int) (int, error)
}

// Client returns a exist client.
func Client() Factory {
	return client
}

// SetClient set client to global client.
func SetClient(factory Factory) {
	client = factory
}
