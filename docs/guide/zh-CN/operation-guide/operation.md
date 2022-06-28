## 启动命令

```yaml
# pigpig 全配置

# 服务配置
server:
  healthz: true # 是否开启健康检查，如果开启会安装 /healthz 路由，默认 true
  middlewares: requestid,logger # 加载的 中间件列表，多个中间件，逗号(,)隔开
  max-ping-count: 3 # http 服务启动后，自检尝试次数，默认 3
  cluster:
    enable: false # 是否开启集群模式
    name: "pigpig-cluster" # 集群名称
    role: "salve"  # 服务器的集群角色
    is_master_handle: true # master 节点是否处理请求, 流量较大时可关闭此选项


# HTTP 配置
insecure:
  bind-address: 0.0.0.0 # 绑定的不安全 IP 地址，设置为 0.0.0.0 表示使用全部网络接口，默认为 127.0.0.1
  bind-port: 8088 # 提供非安全认证的监听端口，默认为 8080
#    bind-address: ${PIGPIG_INSECURE_BIND_ADDRESS} # 绑定的不安全 IP 地址，设置为 0.0.0.0 表示使用全部网络接口，默认为 127.0.0.1
#    bind-port: ${PIGPIG_INSECURE_BIND_PORT} # 提供非安全认证的监听端口，默认为 8080

# HTTPS 配置
secure:
  bind-address: 0.0.0.0 # HTTPS 安全模式的 IP 地址，默认为 0.0.0.0
  bind-port: 8443
  # 使用 HTTPS 安全模式的端口号，设置为 0 表示不启用 HTTPS，默认为 8443
  tls:
    #cert-dir: .iam/cert # TLS 证书所在的目录，默认值为 /var/run/iam
    #pair-name: iam # TLS 私钥对名称，默认 iam
    cert-key:
      cert-file: "./configs/cert/PigPig.crt" # 包含 x509 证书的文件路径，用 HTTPS 认证
      private-key-file: "./configs/cert/PigPig.key" # TLS 私钥

# Redis 配置
redis:
  host: "127.0.0.1" # redis 地址，默认 127.0.0.1:6379
  port: 6379 # redis 端口，默认 6379
  password: "" # redis 密码
  #addrs:
  #master-name: # redis 集群 master 名称
  #username: # redis 登录用户名
  #database: # redis 数据库
  #optimisation-max-idle:  # redis 连接池中的最大空闲连接数
  #optimisation-max-active: # 最大活跃连接数
  #timeout: # 连接 redis 时的超时时间
  #enable-cluster: # 是否开启集群模式
  #use-ssl: # 是否启用 TLS
  #ssl-insecure-skip-verify: # 当连接 redis 时允许使用自签名证书

# etcd 配置
etcd:
  endpoints: "etcd:5566"
  lease-expire: 10
  timeout: 10
  request-timeout: 2
  health_beat_path_prefix: /pigpig_health_beat
  namespace: /pigpig
  UseTLS: false


log:
  name: pigpig # Logger的名字
  development: true # 是否是开发模式。如果是开发模式，会对DPanicLevel进行堆栈跟踪。
  level: debug # 日志级别，优先级从低到高依次为：debug, info, warn, error, dpanic, panic, fatal。
  format: console # 支持的日志输出格式，目前支持console和json两种。console其实就是text格式。
  enable-color: true # 是否开启颜色输出，true:是，false:否
  disable-caller: false # 是否开启 caller，如果开启会在日志中显示调用日志所在的文件、函数和行号
  disable-stacktrace: false # 是否再panic及以上级别禁止打印堆栈信息
  output-paths: ./pigpig.log,stdout # 支持输出到多个输出，逗号分开。支持输出到标准输出（stdout）和文件。
  error-output-paths: ./pigpig.error.log # zap内部(非业务)错误日志输出路径，多个输出，逗号分开

```

### 命令行启动示例
```bash
$ pigpig -c configs/pigpig.yaml --server.cluster.enable=true --server.cluster.name=my_cluster \
 --server.cluster.role=master \
 --log.level=info --redis.host=host:port
```