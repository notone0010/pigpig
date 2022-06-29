// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	v1 "github.com/notone0010/pigpig/internal/pigpig/dudu"
	"github.com/spf13/pflag"

	"github.com/notone0010/pigpig/internal/pkg/server"
)

// ServerRunOptions contains the options while running a generic api server.
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
	Plugins     []string `json:"plugins" mapstructure:"plugins"`

	Cluster ClusterOptions `json:"cluster" mapstructure:"cluster"`
}

// ClusterOptions contains role, enable and so on of current server.
type ClusterOptions struct {
	Enable         bool   `json:"enable" mapstructure:"enable"`
	Role           string `json:"role" mapstructure:"role"`
	IsMasterHandle bool   `json:"is_master_handle" mapstructure:"is_master_handle"`

	// Name is this cluster name
	Name string `json:"name" mapstructure:"name"`
	// LoadPolicy can choose rr, shuffle and consistent-hash
	LoadPolicy string `json:"policy" mapstructure:"policy"`
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig()

	return &ServerRunOptions{
		Healthz:     defaults.Healthz,
		Middlewares: defaults.Middlewares,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares
	c.Cluster = v1.Cluster{
		Enable:         s.Cluster.Enable,
		Role:           s.Cluster.Role,
		IsMasterHandle: s.Cluster.IsMasterHandle,
		LoadPolicy:     s.Cluster.LoadPolicy,
		Name:           s.Cluster.Name,
	}

	return nil
}

// Validate checks validation of ServerRunOptions.
func (s *ServerRunOptions) Validate() []error {
	errors := []error{}

	return errors
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.

	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.")

	fs.BoolVar(&s.Cluster.Enable, "server.cluster.enable", s.Cluster.Enable, ""+
		"Whether enable the server cluster. if you want enable the cluster, you may set this button to true")
	fs.StringVar(&s.Cluster.Name, "server.cluster.name", s.Cluster.Name, ""+
		"While you set 'server.cluster.enable' is true, you must set your cluster name")
	fs.StringVar(&s.Cluster.Role, "server.cluster.role", s.Cluster.Role, ""+
		"While you set 'server.cluster.enable' is true, you can allocate what the current server is role [master|slave]")
	fs.BoolVarP(&s.Cluster.IsMasterHandle, "server.cluster.is-master-handle", "m", s.Cluster.IsMasterHandle, ""+
		"Expression the server whether the work works when the server role is master")
	fs.StringVar(&s.Cluster.LoadPolicy, "server.cluster.policy", s.Cluster.LoadPolicy, ""+
		"The server use what load-balance algorithm when the server is master. If this options is empty default algorithm will be use round-robin")
}
