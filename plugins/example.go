// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/notone0010/pigpig/internal/pigpig/dudu"
	"github.com/notone0010/pigpig/internal/pkg/plugins"
)

type example struct{}

// NewPlugin function name is necessary NewPlugin
func NewPlugin() plugins.PigPigPlugins {
	return &example{}
}

var _ plugins.PigPigPlugins = (*example)(nil)

// ModifyRequest is an optional function that modifies the
// request before send to the remote.
func (e example) ModifyRequest(r *dudu.RequestDetail) {
	log.Printf("modify request")
}

// ModifyResponse is an optional function that modifies the
// Response from the remote. It is called if the remote
// returns a response at all, with any HTTP status code.
// If the backend is unreachable, the optional ErrorHandler is
// called without any call to ModifyResponse.
func (e example) ModifyResponse(r *dudu.RequestDetail, resp *dudu.ResponseDetail) {
	log.Println("modify response")
}

// ModifyError is an optional function if the remote is unreachable then call modify error
func (e example) ModifyError(c *dudu.RequestDetail, errors []error) {
}
