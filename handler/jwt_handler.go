package handler

import (
	"context"
	pb "github.com/MrGaoRock666/Mairuida/user_service/pb"
	"github.com/MrGaoRock666/Mairuida/user_service/middleware"
	"google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
)

// GetUserIDByJWT 是一个测试接口，用于验证 JWT 拦截器是否成功提取 user_id
func (s *UserService) GetUserIDByJWT(ctx context.Context, _ *pb.Empty) (*pb.UserIDResponse, error) {
	// 从 context 中提取 user_id
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "无法从 JWT 中解析 user_id")
	}

	// 成功解析，返回
	return &pb.UserIDResponse{UserId: userID}, nil
}
