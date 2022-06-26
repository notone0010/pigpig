// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// middleware
package middleware

import (
	"net/http/httputil"
	"strings"
	"time"

	v1 "github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/pkg/log"
)

func Recovery() v1.HandlerFunc {
	return func(c *v1.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				headersToStr := strings.Join(headers, "\r\n")

				log.Errorf("[Recovery] %s panic recovered:\n%s\n%s", timeFormat(time.Now()), err, headersToStr)
				// If the connection is dead, we can't write a status to it.
				c.Errors = append(c.Errors, err.(error)) // nolint: errcheck
				c.Abort()
			}
		}()
		c.Next()
	}
}

// timeFormat returns a customized time string for logger.
func timeFormat(t time.Time) string {
	return t.Format("2006/01/02 - 15:04:05")
}
