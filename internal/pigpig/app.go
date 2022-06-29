// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pigpig

import (
	"github.com/notone/pigpig/internal/pigpig/config"
	"github.com/notone/pigpig/internal/pigpig/options"
	"github.com/notone/pigpig/pkg/app"
	"github.com/notone/pigpig/pkg/log"
)

const commandDesc = `The PigPig is a distributed proxy server and it
send your request to what you require remote server, the PigPig support service
automatic register and discover, internal load balance, and plugins
handles the request traffic.

Find more pigpig information at:
    https://github.com/notone0010/pigpig/blob/master/README.md`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("PigPig Proxy Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
