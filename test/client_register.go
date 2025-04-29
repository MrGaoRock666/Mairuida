package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/MrGaoRock666/Mairuida/user_service/pb"
)

func main() {
	// 设定连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 连接 user_service 服务
	conn, err := grpc.DialContext(ctx, "localhost:5001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("无法连接服务: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewUserServiceClient(conn)

	// 构造注册请求（每次都改一个用户名/邮箱/手机号避免重复）
	req := &pb.RegisterRequest{
		Username: "test002",
		Password: "123456",
		Email:    "test002@example.com",
		Phone:    "13800000002",
	}

	// 调用注册接口
	resp, err := client.RegisterUser(context.Background(), req)
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	fmt.Println("注册响应:")
	fmt.Println("成功？", resp.Success)
	fmt.Println("信息：", resp.Message)
	fmt.Println("用户ID:", resp.UserId)
}
