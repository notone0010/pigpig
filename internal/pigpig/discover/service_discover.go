// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package discover

import (
	"context"
)

// ServiceDiscover service discover interface.
type ServiceDiscover interface {
	DiscoverStart(ctx context.Context, prefix string) error

	GetService() []string

	GetHttpsService() []string

	RestartSession() error
	Close() error
}
