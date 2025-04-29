package main

import (
	"log"
	"net"

	"github.com/MrGaoRock666/Mairuida/order_service/config"
	"github.com/MrGaoRock666/Mairuida/order_service/handler"
	"github.com/MrGaoRock666/Mairuida/order_service/pb"

	"google.golang.org/grpc"
)

func main() {
	// 加载运费计算配置
	shippingConfig, err := config.LoadShippingConfig()
	if err != nil {
		log.Fatalf("Failed to load shipping config: %v", err)
	}

	// 初始化数据库
	db := config.InitOrderDB()

	// 初始化 Redis 连接，确保 InitRedis 函数返回正确的 Redis 客户端实例
	redisClient := config.InitRedis()
	if redisClient == nil {
		log.Fatalf("Failed to initialize Redis client")
	}

	// 创建 gRPC 服务
	server := grpc.NewServer()

	// 注册服务，并将数据库连接、Redis 连接和配置注入
	orderService := &handler.OrderService{
		DB:             db,
		RedisClient:    redisClient,
		ShippingConfig: shippingConfig,
	}
	pb.RegisterOrderServiceServer(server, orderService)

	// 监听端口(5002)
	lis, err := net.Listen("tcp", ":5002")
	if err != nil {
		log.Fatalf("Listen failed:%v", err)
	}

	log.Println("OrderService is started on port 5002...")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Start failed:%v", err)
	}
}
