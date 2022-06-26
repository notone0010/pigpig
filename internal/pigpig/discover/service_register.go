// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// discover
package discover

import "context"

type ServiceRegister interface {
	NewServiceRegister(ctx context.Context, cluster, proto, key, value string) error
	RecoveryServiceRegister(ctx context.Context) error

	RestartSession() error
	Close() error
}
