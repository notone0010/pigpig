// // Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logrus

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func NewLogger(zapLogger *zap.Logger) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logger.AddHook(newHook(zapLogger))

	return logger
}
