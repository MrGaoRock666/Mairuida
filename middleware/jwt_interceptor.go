// JWT拦截器 - 用于对所有 gRPC 请求统一进行 JWT 鉴权处理
package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"    // 引入 JWT 库
	"google.golang.org/grpc"          // gRPC 核心库
	"google.golang.org/grpc/codes"    // gRPC 错误码
	"google.golang.org/grpc/metadata" // 用于获取请求的 metadata（头部信息）
	"google.golang.org/grpc/status"   // 用于返回 gRPC 错误状态
)

// 自定义 context 的 key 类型，防止跟其他 context 键冲突
type contextKey string

// 常量定义
const (
	JWT_SECRET       = "your_secret_key"     // JWT 密钥（必须和你生成 token 时使用的密钥一致）
	ContextUserIDKey = contextKey("user_id") // 用于在 context 中存储 user_id 的 key
)

// JWTAuthInterceptor 返回一个 gRPC 的拦截器函数，用于拦截所有请求并进行 JWT 校验
func JWTAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		//context：它是一个上下文对象，包含了请求的生命周期信息，比如请求是否超时、是否取消等。
		//每个 gRPC 调用都需要一个 context，它用于追踪这个调用的状态。
		ctx context.Context, // 上下文对象，传递认证信息、请求状态等
		req interface{}, // 请求体
		info *grpc.UnaryServerInfo, // 当前 RPC 方法的相关信息
		handler grpc.UnaryHandler, // 实际业务处理函数
	) (interface{}, error) {

		// 放行不需要认证的接口（你可以根据 proto 方法名修改）
		// info.FullMethod 形如 "/pb.UserService/LoginUser"
		if strings.HasSuffix(info.FullMethod, "LoginUser") ||
			strings.HasSuffix(info.FullMethod, "RegisterUser") ||
			strings.HasSuffix(info.FullMethod, "LoginByCode") ||
			strings.HasSuffix(info.FullMethod, "SendLoginCode") {
			return handler(ctx, req)
		}
		//metadata：metadata 就是请求的额外信息，类似 HTTP 请求头的概念。
		//在 gRPC 中，我们通常会把 JWT Token 放在 metadata 里，并通过 context 传递给后端服务。
		// 尝试从 gRPC metadata 中获取请求头部信息(从context中获取metadata)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// 获取失败，说明请求头无效，返回未认证错误
			return nil, status.Error(codes.Unauthenticated, "metadata read error")
		}

		// 尝试从 metadata 中读取名为 "authorization" 的字段
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			// 如果 authorization 为空，说明未携带 token
			return nil, status.Error(codes.Unauthenticated, "Authorization head is empty")
		}

		// 通常 Authorization 头会长这样：Bearer xxxxxx
		// 使用 strings.TrimPrefix 去除 "Bearer " 前缀，提取出真正的 token 值
		tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")

		// 如果没有以 "Bearer " 开头，说明格式不对，拦截请求
		if tokenStr == authHeaders[0] {
			return nil, status.Error(codes.Unauthenticated, "format error: lack Bearer begin")
		}

		// 开始解析并验证 token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// 这里返回的是用于签名验证的密钥
			return []byte(JWT_SECRET), nil
		})

		// 如果解析失败或 token 无效（过期、伪造等），拦截请求
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "Token invalid or outtime")
		}

		// 将 token 中的 payload（claims）解析为 MapClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// 如果无法读取 claims（载荷），说明 token 非法
			return nil, status.Error(codes.Unauthenticated, "Token can't read claims")
		}

		// 从 claims 中获取自定义字段 "user_id"，我们在生成 token 时手动加入的
		userID := claims["user_id"].(string)

		// 将 user_id 存储进 context，方便后续业务 handler 中使用
		ctx = context.WithValue(ctx, ContextUserIDKey, userID)

		// 放行，调用真正的 handler 处理请求
		return handler(ctx, req)
	}
}

// 提供一个函数用于从 context 中提取 user_id，供 handler 使用
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	// 尝试从 context 中读取 user_id
	uid, ok := ctx.Value(ContextUserIDKey).(string)
	return uid, ok
}
