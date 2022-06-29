// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/notone0010/pigpig/internal/pkg/server"
)

// ProxyHandlerOptions proxy handler options.
type ProxyHandlerOptions struct {
	// ForceProxyHttps will enable that deal with https stream
	ForceProxyHttps bool

	// HttpServerOptions http/https server
	HttpServerOptions  *server.InsecureServingInfo
	HttpsServerOptions *server.SecureServingInfo

	// Plugins work within proxy server processing
	Plugins []string
}

// NewProxyHandlerOptions creates a NewProxyHandlerOptions object with default parameters.
func (s *ProxyHandlerOptions) NewProxyHandlerOptions() *ProxyHandlerOptions {
	return &ProxyHandlerOptions{
		ForceProxyHttps: false,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ProxyHandlerOptions) ApplyTo(c *server.Config) error {
	s.HttpServerOptions = c.InsecureServing
	s.HttpsServerOptions = c.SecureServing

	return nil
}
