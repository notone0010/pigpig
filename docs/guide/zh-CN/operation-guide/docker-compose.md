## docker-compose.yaml 示例

```yaml
version: '3'

# 指定集群网络信息
networks:
  pigpig_net:
    driver: bridge
    ipam:
      driver: default
      config:
        - subnet: 168.171.10.0/24

# 所有服务均在同一个network中
services:
  pigpig-etcd:
    image: quay.io/coreos/etcd
    container_name: pigpig-etcd
    command: etcd -name etcd -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379 -listen-peer-urls http://0.0.0.0:2380
    ports:
      - "5566:2379"
      - "5567:2380"
    networks:
      pigpig_net:
        ipv4_address: 168.171.10.4 # 绑定IP地址
        aliases:
          - etcd      # 指定服务别名，方便IP变更时可用别名访问
    environment:
      - ETCDCTL_API=3   # 指定etcd api版本
    restart: unless-stopped

  pigpig-redis:
    image: library/redis:5.0
    container_name: pigpig-redis
    command: redis-server --bind 0.0.0.0  # 指定绑定网卡，默认是127.0.0.1 否则无法从外部访问
    networks:
      pigpig_net:
        ipv4_address: 168.171.10.5
        aliases:
          - redis
    restart: unless-stopped

  pigpig-master:
    image: 127.0.0.1:5050/notone/pigpig-arm64:v1.0.0-5-g080201b  # pigpig镜像地址
    container_name: pigpig-master
    # 启动命令 是否开启集群，集群名称等
    command: -c /etc/pigpig/pigpig.yaml --server.cluster.enable=true --server.cluster.name=my_cluster --server.cluster.role=master -m=true --log.level=info
    ports:
      - 8080:8080
      - 8443:8443
    networks:
      pigpig_net:
        ipv4_address: 168.171.10.2
    restart: unless-stopped
    volumes:
      - type: bind
        source: ./configs/pigpig_container.yaml  # 绑定配置文件
        target: /etc/pigpig/pigpig.yaml
        volume:
          nocopy: true
      - type: bind
        source: ./configs/cert
        target: /etc/cert
        volume:
          nocopy: true
    depends_on:
      - pigpig-etcd
      - pigpig-redis

  pigpig-slave1:
    image: 127.0.0.1:5050/notone/pigpig-arm64:v1.0.0-5-g080201b
    container_name: pigpig-slave1
    command: -c /etc/pigpig/pigpig.yaml --server.cluster.enable=true --server.cluster.name=my_cluster --server.cluster.role=slave --log.level=info
    ports:
      - 9080:8080
      - 9443:8443
    networks:
      pigpig_net:
        ipv4_address: 168.171.10.3
    restart: unless-stopped
    volumes:
      - type: bind
        source: ./configs/pigpig_container.yaml
        target: /etc/pigpig/pigpig.yaml
        volume:
          nocopy: true
      - type: bind
        source: ./configs/cert
        target: /etc/cert
        volume:
          nocopy: true
    depends_on:
      - pigpig-etcd
      - pigpig-redis
      - pigpig-master
```