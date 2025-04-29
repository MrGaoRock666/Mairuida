package test

import (
	"context"
	"log" // 用于调试输出
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/MrGaoRock666/Mairuida/user_service/pb" // 替换成你的实际包路径
)

func TestGetUserIDByJWT(t *testing.T) {
	// 1. 连接 gRPC
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// 2. 登录，拿到 JWT Token
	loginResp, err := client.LoginUser(context.Background(), &pb.LoginRequest{
		Username: "your_test_username",
		Password: "your_test_password",
	})
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	// 打印登录返回的 token，确保它不为空
	log.Printf("登录成功，获取到的 Token: %s", loginResp.Token)

	token := loginResp.Token
	if token == "" {
		t.Fatalf("登录成功，但未获取到 Token")
	}

	// 3. 创建 metadata，加入 Authorization Bearer Token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 4. 调用带 JWT 的接口
	resp, err := client.GetUserIDByJWT(ctx, &pb.Empty{})
	if err != nil {
		t.Fatalf("调用失败: %v", err)
	}

	// 打印从 JWT 中提取到的 user_id
	t.Logf("从 JWT 中提取到的 user_id: %s", resp.UserId)
}
