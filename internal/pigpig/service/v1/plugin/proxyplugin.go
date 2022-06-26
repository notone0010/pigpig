// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// proxyplugin
package plugin

import (
	"github.com/notone/pigpig/internal/pigpig/dudu"
)

type ProxyPlugin struct {
	name string
}

// var _ Plugin = (ProxyPlugin)(nil)

func NewProxyPlugin() *ProxyPlugin {
	return &ProxyPlugin{}
}

func (p ProxyPlugin) BeforeSendRequest(c *dudu.Context) {
	// you can do something for before send request in here but you must called c.Next()
	c.RequestDetail.Header.Set("Test-Key", "19960218")
	return
}

func (p ProxyPlugin) BeforeSendResponse(c *dudu.Context) {
	val := c.RequestDetail.Header.Get("Test-Key")
	if val != "" {
		c.ResponseDetail.Header.Set("Test-Key", val)
	}
	return
}
