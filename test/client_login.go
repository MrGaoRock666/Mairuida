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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:5001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("无法连接服务: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// 登录请求（请确保用户名密码存在于数据库）
	req := &pb.LoginRequest{
		Username: "test002",
		Password: "123456",
	}

	resp, err := client.LoginUser(context.Background(), req)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Println("登录成功！")
	fmt.Println("用户ID：", resp.UserId)
	fmt.Println("Token：", resp.Token)
}
