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

// NewPlugin for plugins example.
func NewPlugin() plugins.PigPigPlugins {
	return &example{}
}

// ModifyRequest for plugins example.
func (e example) ModifyRequest(r *dudu.RequestDetail) {
	log.Printf("modify request")
}

// ModifyResponse for plugins example.
func (e example) ModifyResponse(r *dudu.RequestDetail, resp *dudu.ResponseDetail) {
	log.Println("modify response")
}

// ModifyError ...
func (e example) ModifyError(c *dudu.RequestDetail, errors []error) {
}
