// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// pigpig
package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/notone/pigpig/internal/pigpig"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	pigpig.NewApp("pigpig").Run()
}




