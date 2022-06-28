# PigPig - 分布式代理服务器

***PigPig*** 是受**anyproxy**的启发并且结合工作中所遇到的一些问题的解决方案，基于go语言开发的分布式代理服务器。
支持插件化的HTTP/HTTPS代理服务器，支持动态水平扩容。可供二次开发更灵活的处理请求逻辑，并且支持HTTPS的解析

Github地址: https://github.com/notone0010/pigpig

## 安装

1. Docker/Docker-Compose
2. 源码安装

### Pre-requisite (Docker/Docker-Compose)

* Docker 18.03+
* Docker Compose 1.24+

### Pre-requisite (源码安装)

* Go 1.16+
* Redis 5.x+
* Etcd v3

# 快速开始

```bash
$ git clone https://github.com/notone0010/pigpig.git $GOPATH/src/github.com/notone/pigpig
$ cd $GOPATH/src/github.com/notone/pigpig
$ docker-compose up -d
```
命令执行完成后，可以访问 `http://localhost:8080/version` 查看当前服务版本信息

![version](./docs/images/version_small.png)

## 构建
1. git clone 
```bash
$ git clone https://github.com/notone0010/pigpig.git $GOPATH/src/github.com/notone/pigpig
```

2. 编译

```bash
$ cd $GOPATH/src/github.com/notone/pigpig
$ make
```

构建后的二进制文件保存在 `_output/platforms/linux/amd64/` 目录下。

### 使用指南&&参数详解

[Documentation](docs)