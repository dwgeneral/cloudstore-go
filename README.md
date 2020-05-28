# cloudstore-go
基于 Golang 实现的分布式云存储服务

## 项目介绍

本项目的构建初衷是想开发一个私密的家庭相册产品，随着对项目的梳理，感觉自己应该先做一个独立的云存储项目，于是有了本项目的初步规划。在立项，技术选型时，本着边学边练的精神，将自己感兴趣而没有实践过的技术运用一下，岂不美哉。于是选择了Golang，MySQL（对，你没看错，我工作用的是MongoDB），gRPC 等技术栈。

该项目计划实现的功能是：

- 账号体系
  - 用户需要注册、登陆，才能上传文件
  - 用户只能看到自己上传的文件

- 文件上传服务
  - 用户可以自由上传文件，小文件，大文件均可
  - 用户上传文件的流程应该是顺滑的，大文件需要分片上传，断电续传
  - 用户上传一份已经在云端的文件，应该要做到秒传，无论这个文件是谁传上来的
  - 用户上传的文件不能随意丢失，需要连接到云端OSS服务

- 文件下载服务
  - 用户可以下载自己传过的文件

## 技术栈介绍

- Gin
  - 基于 Golang 的一个 Web 应用框架
  - 底层是基于 Go 的 net/http 包
  - 更好的性能和更快的路由
  - 编程体验优秀。只需要引入包、定义路由、编写Handler即可开发应用
  - 简单，理解了核心结构 gin.Context 即可使用 Gin 流畅编程
  - 简洁，没有提供 ORM，CONFIG 等组件，把选择权留给开发者

- MySQL
  - 使用 MySQL 作为持久化数据库
  - 主要存储用户账户信息，用户文件关系，文件元信息等数据

- Redis
  - 使用 Redis 作为缓存数据库
  - 主要缓存在分片上传大文件时的状态信息 uploadID, chunkID, index, filesize 等

- RabbitMQ
  - 使用 RabbitMQ 作为消息队列
  - 将文件异步存储至OSS

- go-micro
  - 基于 Golang 的一个插件式微服务框架
  - 提供服务发现、负载均衡、同步/异步通信、服务接口等，所有组件均为 Interface，便于扩展
  - 服务间传输数据格式为 protobuf，效率高，安全
  - 主要组件
    - Registry 
      - 服务发现、发现、注销、监测机制
      - 服务注册中心支持 consul、etcd、zookeeper、gossip、k8s、eureka等
    - Select
      - 选择器提供了负载均衡，可以通过过滤方法对微服务进行过滤，并通过不同路由算法选择微服务，以及缓存等
    - Transport
      - 微服务间同步请求/响应通信方式，相对Go标准net包做了更高的抽象，支持更多的传输方式，如http、grpc、tcp、udp、Rabbitmq等
    - Broker
      - 微服务间异步发布/订阅通信方式，更好的处理分布式系统解耦问题，默认使用http方式，生产环境通常会使用消息中间件，如Kafka、RabbitMQ、NSQ等
    - Codec
      - 服务间消息的编解码，支持json、protobuf、bson、msgpack等，与普通编码格式不同都是支持RPC格式
    - Server
      - 用于启动服务，为服务命名、注册Handler、添加中间件等
    - Client
      - 提供微服务客户端，通过Registry、Selector、Transport、Broker实现以服务名来查找服务、负载均衡、同步通信、异步消息等

- Docker
  - 容器化是必须的

- Kubernetes
  - 既上次搭建完K8S集群后，好久没接触过了，这次重新熟悉一下

- Istio (TODO)
  - Service Mesh 据说好像很不错

- Kafka (TODO)
  - 如果未来有日志处理需求了的话，可以考虑

- React (TODO)
  - 前端界面很丑，我得用 React Ant Design 美化一下
