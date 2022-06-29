// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/notone0010/pigpig/internal/pigpig/dudu"
	"github.com/notone0010/pigpig/internal/pkg/plugins"
)

type Example struct {
}

func NewPlugin() plugins.PigPigPlugins {
	return &Example{}
}

func (e Example) ModifyRequest(r *dudu.RequestDetail) {
	log.Printf("modify request")
}

func (e Example) ModifyResponse(r *dudu.RequestDetail, resp *dudu.ResponseDetail) {
	log.Println("modify response")

}

func (e Example) ModifyError(c *dudu.RequestDetail, errors []error) {
}
