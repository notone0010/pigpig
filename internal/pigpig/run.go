// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pigpig

import "github.com/notone/pigpig/internal/pigpig/config"

// Run run a server.
func Run(cfg *config.Config) error {
	server, err := createProxyServer(cfg)
	if err != nil {
		return err
	}


	return server.PrepareRun().Run()
}
