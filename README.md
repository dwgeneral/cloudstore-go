# cloudstore-go
基于 Golang 的分布式云存储服务

#### docker 启动 mysql 主从
- 创建 mysql 文件夹，建立 master & slave 配置文件
- 通过 docker 启动两个 mysql 实例，一个作为 master，一个作为 slave
```shell
$ vi mysql/master.conf
$ vi mysql/slave.conf

$ docker pull mysql:5.7

$ docker run -d --name mysql-master -p 13306:3306 -v ~/go/src/cloudstore-go/mysql/master.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v ~/go/src/cloudstore-go/mysql/db:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

$ docker run -d --name mysql-slave -p 13307:3306 -v ~/go/src/cloudstore-go/mysql/slave.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v ~/go/src/cloudstore-go/mysql/db_slave:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
```

- 登陆 Master 节点，进行主从配置
```shell
# 192.168.123.xx 是你本机的内网ip (ifconfig查看)
$ mysql -u root -h 192.168.123.xx -P13306 -p123456
mysql> GRANT REPLICATION SLAVE ON *.* TO 'slave'@'%' IDENTIFIED BY 'slave';
mysql> flush privileges;
mysql> create database cloudstore default character set utf8mb4;
mysql> show master status \G;
*************************** 1. row ***************************
             File: log.000003
         Position: 1313
     Binlog_Do_DB:
 Binlog_Ignore_DB:
Executed_Gtid_Set:
```

- 登陆 Slave 节点，进行主从配置
```shell
$ mysql -u root -h 192.168.123.xx -P13307 -p123456
mysql> stop slave;
# 注意其中的日志文件和数值要和上面show master status的值对应
mysql> CHANGE MASTER TO MASTER_HOST='your ip',master_port=13306,MASTER_USER='slave',MASTER_PASSWORD='slave',MASTER_LOG_FILE='log.log.000003',MASTER_LOG_POS=0;
mysql> start slave;
mysql> show slave status G;
// ...
Slave_IO_Running: Yes 
Slave_SQL_Running: Yes 
// ...
```
- 配置完成，此时在 master 节点的数据与 slave 节点的数据会通过binlog进行同步

#### 使用 MySQL 技术概览
- 通过 sql.DB 来管理数据库连接对象
- 通过 sql.Open 来创建协程安全的 sql.DB 对象
  - 一般来说这个对象是作为长连接来使用的
  - 我们不需要频繁的调用 Open / Close 方法
  - 减少资源消耗和服务器压力
- 优先使用 Prepared Statement
  - 防止SQL注入攻击
  - 比手动拼接字符串更有效
  - 方便实现自定义参数查询

#### docker 启动 RabbitMQ 服务
- 创建数据卷挂载目录 ./db/rabbitmq/
- 通过 docker 启动 rabbitmq 服务，默认 5672 为 rabbitmq 服务端口，15672 为Web管理界面端口，25672 为集群间节点通信端口 
- RabbitMQ Web 管理后台默认用户名及密码都为 guest
```shell
$ docker run -d --hostname rabbit-server --name rabbit -p 5672:5672 -p 15672:15672 -p 25672:25672 -v ~/go/src/cloudstore-go/db/rabbitmq:/var/lib/rabbitmq rabbitmq:management
```
- 启动完成后，你可以登陆到 http://localhost:15672 尝试添加 exchange，queue，publish message, get message 体验一下 RabbitMQ 基本的消息流转逻辑, 如果遇到问题，请参考：https://www.rabbitmq.com/getstarted.html 官方文档


#### 安装 go-micro、Protobuf、consul、grpc 等工具集
- 安装教程参考：https://studygolang.com/articles/16832

- 访问服务发现Consul管理后台 http://localhost:8500

#### 微服务架构迁移
- API网关
  - 将用户的HTTPS请求更换为微服务内部通信协议
  - 限流、熔断等功能

- 配置中心

- 服务发现(服务注册中心) Consul
  - 提供服务发现功能

- 用户账户服务 UserService
  - 提供用户登陆、注册功能

- 文件下载服务 DownloadService
  - 与上传服务分开，解耦，高并发时相互之间不再产生影响

- 文件上传服务 UploadService
  - 文件直传功能
  - 大文件分片上传 / 秒传 / 断点续传 功能

- 文件转移服务 TransferService
  - 将上传的文件异步发送给OSS云存储

- 消息队列 RabbitMQ

- DBProxy
  - 将数据层抽象，与具体的数据库实现解耦
  - MySQL
  - Redis

#### 项目清单
- 功能列表
  - 文件上传服务
    - 文件下载
    - 秒传功能 / 断点续传
  - 文件转移服务
    - Go结合RabbitMQ消息队列实现文件异步存储到OSS
  - 账号系统

- 技术栈列表
  - Go
  - Gin
  - MySQL
  - Redis
  - Docker
  - RabbitMQ
  - go-micro
  - gPRC
  - Kubernetes
