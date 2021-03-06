#!/usr/bin/env bash

# Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# Common utilities, variables and checks for all build scripts.
set -o errexit
set +o nounset
set -o pipefail

# Sourced flag
COMMON_SOURCED=true

# The root of the build/dist directory
PIGPIG_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${PIGPIG_ROOT}/scripts/lib/init.sh"
source "${PIGPIG_ROOT}/scripts/install/environment.sh"

# 不输入密码执行需要root权限的命令
#function pigpig::common::sudo {
#  echo ${LINUX_PASSWORD} | sudo -S $1
#}
