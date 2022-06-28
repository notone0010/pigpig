// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dudu

import (
	"fmt"
	"net/http"
	"time"
)

// Define colors.
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter.
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes.
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Host is a host the client requests for remote server
	Host string
	// Path is a path the client requests.
	Path string
	// Protocol is a HTTP protocol given to the request
	Protocol string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return Green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return White
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return Yellow
	default:
		return Red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return Blue
	case http.MethodPost:
		return Cyan
	case http.MethodPut:
		return Yellow
	case http.MethodDelete:
		return Red
	case http.MethodPatch:
		return Green
	case http.MethodHead:
		return Magenta
	case http.MethodOptions:
		return White
	default:
		return Reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return Reset
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	// if param.IsOutputColor() {
	// 	statusColor = param.StatusCodeColor()
	// 	methodColor = param.MethodColor()
	// 	resetColor = param.ResetColor()
	// }

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
		// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.ClientIP,
		param.Latency,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
