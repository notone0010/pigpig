// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// pigpig
package pigpig

import (
	"testing"
	"github.com/notone/pigpig/internal/pigpig/options"
	"github.com/notone/pigpig/internal/pigpig/config"
	"github.com/notone/pigpig/pkg/log"
)

func TestRun(t *testing.T) {
	opts := options.NewOptions()
	cfg, _ := config.CreateConfigFromOptions(opts)
	if err := Run(cfg); err != nil{
		log.Fatalf("be starting with error: %s", err.Error())
	}
}
