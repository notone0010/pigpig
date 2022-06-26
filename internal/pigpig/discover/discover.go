// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package discover

//go:generate mockgen -self_package=github.com/marmotedu/iam/internal/apiserver/store -destination mock_store.go -package store github.com/marmotedu/iam/internal/apiserver/store Factory,UserStore,SecretStore,PolicyStore

var client Factory

// Factory defines the iam platform storage interface.
type Factory interface {
	Register() ServiceRegister
	Discover() ServiceDiscover

	Info() InfoStore
	// PolicyAudits() PolicyAuditStore
	Close() error
}

// Client return the discover client instance.
func Client() Factory {
	return client
}

// SetClient set the iam discover client.
func SetClient(factory Factory) {
	client = factory
}
