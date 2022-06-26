// // Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import "github.com/notone/pigpig/pkg/rollinglog"

func main() {
	opts := &rollinglog.Options{
		Level:            "debug",
		Format:           "console",
		EnableColor:      false,
		DisableCaller:    true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
		Rolling:          true,
		RollingMaxSize:   1,
	}
	// 初始化全局logger
	rollinglog.Init(opts)
	defer rollinglog.Flush()

	for i := 0; i < 10000; i++ {
		// rollinglog.Debug("This is a debug message")
		// rollinglog.Warnf("This is a formatted %s message", "hello")
		rollinglog.V(rollinglog.InfoLevel).Info("nice to meet you.")
	}
}
