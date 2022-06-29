// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/notone0010/pigpig/internal/pigpig/dudu"
)

// SendFinalResponse send final response to client.
func (p *ProxyController) SendFinalResponse(c *dudu.Context) error {
	resHeader := c.ResponseDetail.Header
	resBody := c.ResponseDetail.Body

	transferEncoding := resHeader.Get("Transfer-Encoding")
	contentLength := resHeader.Get("Content-Length")
	if transferEncoding != "" {
		resHeader.Set("X-Pigpig-Origin-Transfer-Encoding", transferEncoding)
		resHeader.Del("Transfer-Encoding")
	}
	if contentLength != "" {
		resHeader.Del("Content-Length")
	}
	resHeader.Set("Content-Length", strconv.Itoa(len(resBody)))
	delete(resHeader, "Retry-Count")

	h2 := resHeader.Clone()

	headerSeq := ","
	for k, v := range h2 {
		if k == "Cookie" {
			headerSeq = ";"
		}
		c.Writer.Header().Set(k, strings.Join(v, headerSeq))
	}
	setCookies := c.ResponseDetail.Cookies()
	if len(setCookies) != 0 {
		for _, setCookie := range setCookies {
			http.SetCookie(c.Writer, setCookie)
		}
	}

	// 限速或限流可以在这里操作

	c.Writer.WriteHeader(c.ResponseDetail.StatusCode)
	if _, err := c.Writer.Write(resBody); err != nil {
		return err
	}

	return nil
}
