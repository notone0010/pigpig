// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plugins

import (
	"plugin"
	"strings"

	"github.com/notone0010/pigpig/internal/pigpig/dudu"
	"github.com/notone0010/pigpig/pkg/log"
	"github.com/notone0010/pigpig/pkg/util/fileutil"
)

var defaultPluginsDir = fileutil.GetProjectDirPath() + "/plugin.md"

// PigPigPlugins plugin of the PigPig server interface.
type PigPigPlugins interface {

	// ModifyRequest is an optional function that modifies the
	// request before send to the remote.
	ModifyRequest(c *dudu.RequestDetail)

	// ModifyResponse is an optional function that modifies the
	// Response from the remote. It is called if the remote
	// returns a response at all, with any HTTP status code.
	// If the backend is unreachable, the optional ErrorHandler is
	// called without any call to ModifyResponse.
	ModifyResponse(c *dudu.RequestDetail, r *dudu.ResponseDetail)

	// ModifyError is an optional function if the remote is unreachable then call modify error
	ModifyError(c *dudu.RequestDetail, errors []error)
}

// PluginsOptions ...
type PluginsOptions struct {
	path []string

	Plugins dudu.HandlersChain
}

// NewPluginsOptions returns new plugin.md options.
func NewPluginsOptions(path []string) *PluginsOptions {
	return &PluginsOptions{path: path}
}

type pluginFunc func() PigPigPlugins

// LoadPlugins load plugin.md.
func (o *PluginsOptions) LoadPlugins() {
	log.Infof("attempt to load default plugin.md")
	defaultList := fileutil.ListFile(defaultPluginsDir)
	o.loadPlugins(defaultList)
	log.Infof("loaded default plugin.md")

	if len(o.path) > 0 {
		log.Infof("begin load appoint plugin size -> %d", len(o.path))
		o.loadPlugins(o.path)
	}
}

func (o *PluginsOptions) loadPlugins(f []string) {
	for _, pluginFile := range f {
		if strings.HasSuffix(pluginFile, ".so") {
			p, err := plugin.Open(pluginFile)
			if err != nil {
				log.Warnf("attempt to load %s but encounter an error %s so that skip it", pluginFile, err.Error())
			} else {
				o.pluginComplete(pluginFile, p)
			}
		}
	}
}

func (o *PluginsOptions) pluginComplete(pluginFile string, p *plugin.Plugin) {
	newPlugin, err := p.Lookup("NewPlugin")
	if err != nil {
		log.Fatalf("failed lookup 'NewPlugin' in %s", pluginFile)
	}
	loadPlugin := newPlugin.(func() PigPigPlugins)()
	handleFunc := GetDuDuHandlerFunc(loadPlugin)
	o.Plugins = append(o.Plugins, handleFunc)
	log.Infof("loaded plugin -> %s success", pluginFile)
}

// GetDuDuHandlerFunc compose to dudu.HandlerFUnc.
func GetDuDuHandlerFunc(p PigPigPlugins) dudu.HandlerFunc {
	return func(c *dudu.Context) {
		p.ModifyRequest(c.RequestDetail)
		c.Next()
		if len(c.Errors) > 0 {
			p.ModifyError(c.RequestDetail, c.Errors)
			return
		}
		p.ModifyResponse(c.RequestDetail, c.ResponseDetail)
	}
}
