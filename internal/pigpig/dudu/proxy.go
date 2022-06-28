// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dudu

// Role ...
type Role string

// HandlerFunc defines the handler used by middleware as return value.
type HandlerFunc func(c *Context)

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc

// HandlerErrorFunc defines the handler used by error as return response.
type HandlerErrorFunc func(c *Context, errors []error)

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
