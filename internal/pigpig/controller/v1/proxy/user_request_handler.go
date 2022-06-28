// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"github.com/marmotedu/errors"
	"github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/internal/pkg/code"
	"github.com/notone/pigpig/pkg/core"
	"github.com/notone/pigpig/pkg/log"
)

// UserRequestHandler deal with any http requests.
func (p *ProxyController) UserRequestHandler(c *dudu.Context) {
	c.Handlers = p.Plugins
	c.Handlers = append(c.Handlers, p.srv.Proxy().FetchRemoteResponse)
	c.Next()

	if c.Errors != nil && len(c.Errors) > 0 && p.handleErrorFunc != nil {
		p.handleErrorFunc(c, c.Errors)
	}
	if c.Errors != nil && len(c.Errors) > 0 && c.ResponseDetail == nil {
		aggErr := errors.NewAggregate(c.Errors)
		aggError := errors.WithCode(code.ErrRemoteUnreached, aggErr.Error())
		log.Warnf("one errors: %s", aggErr.Error())
		core.WriteResponse(c.Writer, c.Request, aggError, nil)
		return
	}

	if err := p.SendFinalResponse(c); err != nil {
		log.Errorf("be seem had a error when send final response ----> %s", err.Error())
	}
}
