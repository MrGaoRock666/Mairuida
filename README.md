# 迈睿达--物流管理系统

该项目是一个基于微服务架构的物流管理系统，主要包含 `User Service` 和 `Order Service` 两个核心模块。`User Service` 负责处理与用户相关的业务逻辑，`Order Service` 负责处理订单相关的业务逻辑，同时系统集成了 Kafka 实现异步消息处理，提升系统的响应速度和可扩展性。

## 项目结构
```bash
/Mairuida
├── go.mod                        # Go 模块管理文件
├── go.sum                        # Go 模块依赖清单
├── /proto                        # Protobuf 定义文件（Google 官方及自定义）
│   ├── google/
│   │   └── protobuf/
│   │       ├── timestamp.proto   # 示例 Protobuf 文件
│   ├── user.proto                # 自定义 User 服务的 Proto 文件
│   └── order.proto               # 自定义 Order 服务的 Proto 文件
├── /user_service                 # User Service 实现
│   ├── cmd/
│   │   └── main.go               # 服务启动入口
│   ├── handler/                  # 业务逻辑处理
│   ├── config/                   # 配置信息（如数据库、缓存等）
│   ├── pb/                       # 生成的 Proto 文件代码（user.pb.go）
│   └── util/                     # 工具类，如 Kafka 工具
├── /order_service                # Order Service 实现
│   ├── cmd/
│   │   └── main.go               # 服务启动入口
│   ├── handler/                  # 业务逻辑处理
│   ├── config/                   # 配置信息（如数据库、缓存等）
│   ├── pb/                       # 生成的 Proto 文件代码（order.pb.go）
│   └── util/                     # 工具类
├── /kafka_consumer               # Kafka 消费者服务
│   └── main.go                   # 消费者服务启动入口
└── ...

项目功能
User Service
用户管理
注册：用户可以注册新账号，支持邮箱、手机号等。注册成功后，通过 Kafka 异步发送欢迎邮件通知。
登录：支持传统密码登录及验证码登录。
获取用户信息：通过用户 ID 获取用户详细信息。
更新用户信息：用户可以更新自己的基本信息及地址簿。
用户注销：支持用户注销账户，账户信息会被标记为“已删除”。
用户恢复：支持恢复注销用户的功能。
地址簿管理 用户可以管理多个地址，包括添加、删除和更新地址。
Order Service
订单操作
创建订单：用户可以创建新订单，系统会生成唯一的订单号，使用事务保证数据的一致性。
查询订单详情：根据订单号查询完整的订单信息，先从 Redis 缓存查询，若未命中再从数据库查询并更新缓存。
运费估算：根据订单的距离、重量、体积等信息估算运费，先从 Redis 缓存查询结果，若未命中再进行计算，并将结果存入缓存。
异步消息处理
集成 Kafka 实现用户注册后的异步通知，如发送欢迎邮件，提高系统响应速度。

技术栈
通信协议：使用 gRPC 进行服务间通信，借助 Protobuf 定义接口，生成强类型代码，提高系统的健壮性和开发效率。
数据库：使用 MySQL 存储用户信息和订单数据，通过 GORM 进行数据库操作，在关键操作中使用事务保证数据一致性。
缓存：使用 Redis 缓存热点数据，如订单信息和运费估算结果，减少数据库压力，提升系统响应速度。
消息队列：集成 Kafka 实现异步消息处理，例如用户注册成功后发送欢迎邮件的任务可通过 Kafka 异步处理。

启动与测试
1. 克隆项目

bash
git clone https://github.com/你的账户/Mairuida.git
cd Mairuida
2. 安装依赖
确保已安装 Go 环境，并且在项目根目录执行以下命令：


bash
go mod tidy
3. 启动 ZooKeeper 和 Kafka
如果你使用的是 Kafka 依赖 ZooKeeper 的模式，需要先启动 ZooKeeper 和 Kafka 服务：


bash
# 启动 ZooKeeper（假设 ZooKeeper 安装在 /root/Downloads/apache-zookeeper-3.8.4-bin）
cd /root/Downloads/apache-zookeeper-3.8.4-bin
bin/zkServer.sh start

# 启动 Kafka（假设 Kafka 安装在 /root/Downloads/kafka_2.13-3.6.1）
cd /root/Downloads/kafka_2.13-3.6.1
bin/kafka-server-start.sh config/server.properties
4. 启动 User Service
在 /user_service/cmd 目录下运行：


bash
go run main.go
默认监听端口：5001

5. 启动 Order Service
在 /order_service/cmd 目录下运行：


bash
go run main.go
默认监听端口：5002

6. 启动 Kafka 消费者服务
在 /kafka_consumer 目录下运行：


bash
go run main.go
7. 重新编译 Protobuf 文件
如果你修改了 .proto 文件，需要重新生成对应的 .pb.go 文件。例如，重新生成 user.proto 和 order.proto 文件对应的代码：


bash
# 生成 user.proto 对应的代码
cd /Mairuida/proto
protoc --go_out=../user_service/pb --go-grpc_out=../user_service/pb user.proto

# 生成 order.proto 对应的代码
protoc --go_out=../order_service/pb --go-grpc_out=../order_service/pb order.proto
8. 测试接口
你可以使用 Postman 或其他 gRPC 客户端来调用以下接口：

User Service
注册：POST /user/register
登录：POST /user/login
获取用户信息：GET /user/{id}
更新用户信息：PUT /user/{id}
地址簿管理：GET /user/{id}/addresses
Order Service
创建订单：POST /order/create
查询订单详情：GET /order/{order_id}
运费估算：POST /order/estimate_cost
依赖
数据库 MySQL：用于存储用户信息及订单相关数据。
Redis 用于缓存热点数据，减少数据库压力。
Kafka 用于实现异步消息处理，如用户注册时发送欢迎邮件等。
未来计划
时效预测：结合机器学习预测快递送达时间。
智能调度：基于交通、天气等数据动态调整运输路线。
异步处理：进一步完善 Kafka 和 Redis 的使用，提升系统性能和扩展性。
移动端支持：开发适配移动端的 API 接口。
注意事项
请确保你的 Go 环境已安装。
如果重新生成 .proto 文件，记得执行 go mod tidy 以确保依赖更新。
在生产环境中，建议使用 Docker 或其他容器化方式进行部署。
启动服务前，确保 ZooKeeper 和 Kafka 服务正常运行。
贡献
如果你有任何建议或问题，欢迎提交 issues 或 pull requests。