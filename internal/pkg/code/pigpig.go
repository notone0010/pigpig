// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int

// pigpig: service errors.
const (
	// ErrUserNotFound - 502: Error occurred while requesting remote server.
	ErrRemoteUnreached int = iota + 110201

	// ErrRemoteTimeout - 504: Timeout occurred while requesting remote server.
	ErrRemoteTimeout
	// ErrUserAlreadyExist - 500: Error occurred while remote response content decoding.
	ErrContentDecodingFailed

	// ErrContentEncodingFailed - 500: Error occurred while remote response content encoding.
	ErrContentEncodingFailed

	// ErrProxyInternal - 500: An internal error occurred while the proxy server is processing.
	ErrProxyInternal
)
