// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// proxy
package proxy

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/pkg/log"
)

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
	if setCookies != nil && len(setCookies) != 0 {
		var ()
		for _, setCookie := range setCookies {
			http.SetCookie(c.Writer, setCookie)
		}
	}

	// 限速或限流可以在这里操作

	c.Writer.WriteHeader(c.ResponseDetail.StatusCode)
	code, err := c.Writer.Write(c.ResponseDetail.Body)
	if err != nil {
		log.Errorf("be seem had a error when send final response ----> code: %d", code)
		return err
	}
	return nil
}
