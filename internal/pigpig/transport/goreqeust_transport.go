// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transport

import (
	"net/http"

	"github.com/notone/pigpig/internal/pigpig/dudu"
)

// GoRequestTransport transport interface.
type GoRequestTransport interface {
	// go request transport include functions
	FetchRemoteResponse(c *dudu.RequestDetail) (*http.Response, []byte, error)
}

/*
1.需要处理header
2.处理request body
3.向远端发起请求
4.接收响应并处理
*/
