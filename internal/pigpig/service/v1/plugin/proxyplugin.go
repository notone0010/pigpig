// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plugin

import (
	"github.com/notone0010/pigpig/internal/pigpig/dudu"
)

// ProxyPlugin struct.
type ProxyPlugin struct {
	// name string
}

// var _ Plugin = (ProxyPlugin)(nil)

// NewProxyPlugin return new ProxyPlugin.
func NewProxyPlugin() *ProxyPlugin {
	return &ProxyPlugin{}
}

// BeforeSendRequest sample.
func (p ProxyPlugin) BeforeSendRequest(c *dudu.Context) {
	// you can do something for before send request in here but you must called c.Next()
	c.RequestDetail.Header.Set("Test-Key", "19960218")
}

// BeforeSendResponse sample.
func (p ProxyPlugin) BeforeSendResponse(c *dudu.Context) {
	val := c.RequestDetail.Header.Get("Test-Key")
	if val != "" {
		c.ResponseDetail.Header.Set("Test-Key", val)
	}
}
