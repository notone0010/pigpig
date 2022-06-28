#!/usr/bin/env bash

# Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/iamrnetes/iamrnetes/issues/52255
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
PIGPIG_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${PIGPIG_ROOT}/scripts/lib/util.sh"
source "${PIGPIG_ROOT}/scripts/lib/logging.sh"
source "${PIGPIG_ROOT}/scripts/lib/color.sh"

iam::log::install_errexit

source "${PIGPIG_ROOT}/scripts/lib/version.sh"
#source "${PIGPIG_ROOT}/scripts/lib/golang.sh"
