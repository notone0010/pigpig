// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"fmt"
	"time"

	"github.com/marmotedu/errors"
	v1 "github.com/notone/pigpig/internal/pigpig/dudu"

	"github.com/notone/pigpig/pkg/log"
)

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param v1.LogFormatterParams) string {
	var statusColor, resetColor string

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	return fmt.Sprintf("%s%3d%s - [%s] \"%v %s %s %s %s\" %s",
		// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.ClientIP,
		param.Latency,
		param.Protocol,
		param.Method,
		param.Host,
		param.Path,
		param.ErrorMessage,
	)
}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() v1.HandlerFunc {
	return LoggerWithConfig()
}

//
// // LoggerWithFormatter instance a Logger middleware with the specified log format function.
// func LoggerWithFormatter(f gin.LogFormatter) dudu.HandlerFunc {
// 	return LoggerWithConfig(gin.LoggerConfig{
// 		Formatter: f,
// 	})
// }
//
// // LoggerWithWriter instance a Logger middleware with the specified writer buffer.
// // Example: os.Stdout, a file opened in write mode, a socket...
// func LoggerWithWriter(out io.Writer, notlogged ...string) dudu.HandlerFunc {
// 	return LoggerWithConfig(gin.LoggerConfig{
// 		Output:    out,
// 		SkipPaths: notlogged,
// 	})
// }

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig() v1.HandlerFunc {
	formatter := defaultLogFormatter

	return func(c *v1.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		param := v1.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.Request.RemoteAddr
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.Host = c.Request.Host
		if c.Errors != nil && len(c.Errors) > 0 {
			param.ErrorMessage = errors.NewAggregate(c.Errors).Error()
		}
		if c.Request.TLS != nil {
			param.Protocol = "HTTPS"
		} else {
			param.Protocol = "HTTP"
		}

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		log.L(c).Info(formatter(param))
	}
}
