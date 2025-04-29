// 启动服务
package main

//在main函数中执行自动建表
import (
	"log"
	"net"

	"github.com/MrGaoRock666/Mairuida/user_service/config"
	"github.com/MrGaoRock666/Mairuida/user_service/handler"
	"github.com/MrGaoRock666/Mairuida/user_service/pb"
	"github.com/MrGaoRock666/Mairuida/user_service/util"
	"google.golang.org/grpc"
)

func main() {
	// 初始化数据库
	db := config.InitDB()

	// 初始化 Redis 连接，确保 InitRedis 函数返回正确的 Redis 客户端实例
	redisClient := config.InitRedis()
	if redisClient == nil {
		log.Fatalf("Failed to initialize Redis client")
	}

	// 初始化 Kafka
	util.InitKafka([]string{"localhost:9092"})
	defer util.CloseKafka()

	// 创建 gRPC 服务
	server := grpc.NewServer()

	// 注册服务，并将数据库连接注入
	userService := &handler.UserService{
		DB: db,
	}
	pb.RegisterUserServiceServer(server, userService)

	// 监听端口(5001)
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Listen failed:%v", err)
	}

	log.Println("UserService is started on port 5001...")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Start failed:%v", err)
	}
}
