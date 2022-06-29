// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/notone0010/pigpig/internal/pigpig/transport"
)

// Service defines functions used to return resource interface.
type Service interface {
	Proxy() ProxySrv
}

type service struct {
	transport transport.Factory
}

// NewService returns Service interface.
func NewService(transport transport.Factory) Service {
	return &service{
		transport: transport,
	}
}

func (s *service) Proxy() ProxySrv {
	return newProxy(s)
}
