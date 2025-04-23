物流管理系统 - User Service
该项目是一个基于微服务架构的物流管理系统中的 User Service模块，负责处理与用户相关的业务逻辑，包括用户注册、登录、信息更新、地址簿管理等功能。

项目结构
bash
复制
编辑
/Mairuida
├── go.mod                        # Go模块管理文件
├── go.sum                        # Go模块依赖清单
├── /proto                        # Protobuf定义文件（Google官方及自定义）
│   ├── google/
│   │   └── protobuf/
│   │       ├── timestamp.proto   # 示例Protobuf文件
│   ├── user.proto                # 自定义User服务的Proto文件
├── /user_service                 # User Service实现
│   ├── main.go                   # 服务启动入口
│   ├── handler/                  # 业务逻辑处理
│   ├── config/                   # 配置信息（如数据库、缓存等）
│   ├── pb/                       # 生成的Proto文件代码（user.pb.go）
│   └── ...
项目功能
1. 用户管理
注册：用户可以注册新账号，支持邮箱、手机号等。

登录：支持传统密码登录及验证码登录。

获取用户信息：通过用户ID获取用户详细信息。

更新用户信息：用户可以更新自己的基本信息及地址簿。

用户注销：支持用户注销账户，账户信息会被标记为“已删除”。

用户恢复：支持恢复注销用户的功能。

2. 地址簿管理
用户可以管理多个地址，包括添加、删除和更新地址。

启动与测试
1. 克隆项目
bash
复制
编辑
git clone https://github.com/你的账户/Mairuida.git
cd Mairuida
2. 安装依赖
确保已安装 Go 环境，并且在项目根目录执行以下命令：

bash
复制
编辑
go mod tidy
3. 启动 User Service
在 /user_service 目录下运行：

bash
复制
编辑
go run main.go
默认监听端口：5001

4. 重新编译 Protobuf 文件
如果你修改了 .proto 文件，需要重新生成 user.pb.go 文件。执行以下命令：

bash
复制
编辑
cd /Mairuida/proto
protoc --go_out=../user_service/pb --go-grpc_out=../user_service/pb user.proto
这会将生成的文件放入 /user_service/pb 目录。

5. 测试接口
你可以使用 Postman 或其他 gRPC 客户端来调用以下接口：

注册：POST /register

登录：POST /login

获取用户信息：GET /user/{id}

更新用户信息：PUT /user/{id}

地址簿管理：GET /user/{id}/addresses

6. 数据库配置
目前支持 MySQL 数据库，请确保你的 MySQL 数据库已安装并运行，并且修改 config/database.go 文件中的数据库连接配置。

依赖
1. 数据库
MySQL：用于存储用户信息及相关数据。

2. Redis
用于缓存热点数据，减少数据库压力。

3. Kafka（计划）
用于实现异步消息处理，如用户注册时发送欢迎邮件等。

未来计划
时效预测：结合机器学习预测快递送达时间。

智能调度：基于交通、天气等数据动态调整运输路线。

异步处理：使用 Kafka 和 Redis 提升系统性能和扩展性。

移动端支持：开发适配移动端的 API 接口。

注意事项
请确保你的 Go 环境已安装。

如果重新生成 user.proto 文件，记得执行 go mod tidy 以确保依赖更新。

在生产环境中，建议使用 Docker 或其他容器化方式进行部署。

贡献
如果你有任何建议或问题，欢迎提交 issues 或 pull requests。