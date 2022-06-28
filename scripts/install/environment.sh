#!/usr/bin/env bash

# Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.



# IAM 项目源码根目录
PIGPIG_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# 生成文件存放目录
LOCAL_OUTPUT_ROOT="${IAM_ROOT}/${OUT_DIR:-_output}"

# 设置统一的密码，方便记忆
#readonly PASSWORD=${PASSWORD:-'iam59!z$'}

# Linux系统 going 用户
readonly LINUX_USERNAME=${LINUX_USERNAME:-going}
# Linux root & going 用户密码
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-${PASSWORD}}

# 设置安装目录
readonly INSTALL_DIR=${INSTALL_DIR:-/tmp/installation}
mkdir -p ${INSTALL_DIR}
readonly ENV_FILE=${IAM_ROOT}/scripts/install/environment.sh


# Redis 配置信息
readonly REDIS_HOST=${REDIS_HOST:-127.0.0.1} # Redis 主机地址
readonly REDIS_PORT=${REDIS_PORT:-6379} # Redis 监听端口
readonly REDIS_USERNAME=${REDIS_USERNAME:-''} # Redis 用户名
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-''} # Redis 密码


# pigpig 配置
readonly PIGPIG_DATA_DIR=${PIGPIG_DATA_DIR:-/data/pigpig} # pigpig 各组件数据目录
readonly PIGPIG_INSTALL_DIR=${PIGPIG_INSTALL_DIR:-/opt/pigpig} # pigpig 安装文件存放目录
readonly PIGPIG_CONFIG_DIR=${PIGPIG_CONFIG_DIR:-/etc/pigpig} # pigpig 配置文件存放目录
readonly PIGPIG_LOG_DIR=${PIGPIG_LOG_DIR:-/var/log/pigpig} # pigpig 日志文件存放目录
readonly CA_FILE=${CA_FILE:-${PIGPIG_CONFIG_DIR}/cert/ca.pem} # CA

# server 配置
readonly PIGPIG_HOST=${PIGPIG_HOST:-127.0.0.1} #  部署机器 IP 地址
readonly PIGPIG_INSECURE_BIND_ADDRESS=${PIGPIG_INSECURE_BIND_ADDRESS:-0.0.0.0}
readonly PIGPIG_INSECURE_BIND_PORT=${PIGPIG_INSECURE_BIND_PORT:-8088}
readonly PIGPIG_SECURE_BIND_ADDRESS=${PIGPIG_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly PIGPIG_SECURE_BIND_PORT=${PIGPIG_SECURE_BIND_PORT:-8443}
readonly PIGPIG_SECURE_TLS_CERT_KEY_CERT_FILE=${PIGPIG_SECURE_TLS_CERT_KEY_CERT_FILE:-${PIGPIG_CONFIG_DIR}/cert/PigPig.crt}
readonly PIGPIG_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${PIGPIG_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${PIGPIG_CONFIG_DIR}/cert/PigPig.key}

