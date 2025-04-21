// 启动服务
package main

//在main函数中执行自动建表
import (
	"Mairuida/user_service/config"
	"Mairuida/user_service/handler"
	"Mairuida/user_service/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	//初始化数据库并建表
	db := config.InitDB()

	// 初始化 Redis
	config.InitRedis()

	//创建gRPC服务
	server := grpc.NewServer()

	//注册服务，并将数据库连接注入
	pb.RegisterUserServiceServer(server, &handler.UserService{DB: db})

	//监听端口(5001)
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Listen failed:%v", err)
	}

	log.Println("UserService is started on port 5001...")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Start failed:%v", err)
	}
}
