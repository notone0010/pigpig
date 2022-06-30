// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transport

var client Factory

// Factory defines the iam platform storage interface.
type Factory interface {
	GoRequest() GoRequestTransport

	// you can add new request transport if you need
	Close() error
}

// Client return the transport client instance.
func Client() Factory {
	return client
}

// SetClient set the iam transport client.
func SetClient(factory Factory) {
	client = factory
}
