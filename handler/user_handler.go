// 实现各RPC服务
package handler

import (
	"Mairuida/user_service/model"
	"Mairuida/user_service/config"
	"Mairuida/user_service/pb"
	"context"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9" 
)

// 定义UserServiceServer结构体，实现pb.UserService接口
type UserService struct {
	DB                                *gorm.DB //注入数据库连接
	pb.UnimplementedUserServiceServer          //为没有实现的 RPC 方法提供默认的“未实现”报错处理，从而避免服务 crash
}

// 实现用户注册逻辑
func (s *UserService) RegisterUser(ct context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//使用事务来保证向数据库插入数据过程的原子性
	//开启事务
	tx := s.DB.Begin()
	//创建用户模型实例
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	//查询用户名是否已经存在
	if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err == nil {
		// 用户名已存在
		return &pb.RegisterResponse{
			Success: false,
			Message: "Username already exists!",
		}, nil
	}

	// 查询邮箱是否已存在
	if err := tx.Where("email = ?", req.Email).First(&user).Error; err == nil {
		tx.Rollback()
		return &pb.RegisterResponse{
			Success: false,
			Message: "Email already exists!",
		}, nil
	}

	// 查询电话号码是否已存在
	if err := tx.Where("phone = ?", req.Phone).First(&user).Error; err == nil {
		tx.Rollback()
		return &pb.RegisterResponse{
			Success: false,
			Message: "Phone already exists!",
		}, nil
	}

	//插入数据库(注册失败情况)
	if err := s.DB.Create(&user).Error; err != nil {
		log.Println("Register failed:", err)
		//返回失败响应
		return &pb.RegisterResponse{
			Success: false,
			Message: "Register failed:" + err.Error(),
		}, nil
	}

	//结束事务
	tx.Commit()

	//返回成功情况
	return &pb.RegisterResponse{
		Success: true,
		Message: "Register success!",
		UserId:  fmt.Sprintf("%d", user.ID),
	}, nil
}

// 实现用户登录逻辑
func (s *UserService) LoginUser(ct context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//查询用户是否存在
	var user model.User //GORM 模型中定义的用户表结构

	//如果用户名或密码错误，登录失败
	if err := s.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "username/password error")
	}

	// 模拟生成JWT token（后面再实现真正的JWT）
	token := "fake-jwt-token-" + fmt.Sprintf("%d", user.ID)

	// 返回登录响应
	return &pb.LoginResponse{
		Token:  token,
		UserId: fmt.Sprintf("%d", user.ID),
	}, nil
}

// 实现获取用户信息逻辑（含 Redis 缓存）
func (s *UserService) GetUserInfo(ct context.Context, req *pb.UserIDRequest) (*pb.UserInfo, error) {
	// 构建 Redis 的 key
	cacheKey := fmt.Sprintf("user_info:%s", req.UserId)

	//尝试从 Redis 获取缓存
	cached, err := config.RedisClient.Get(ct, cacheKey).Result()
	if err == nil {
		// Redis 中有数据
		var userInfo pb.UserInfo
		if err := json.Unmarshal([]byte(cached), &userInfo); err == nil {
			log.Println("GetUserInfo: return Redis cachedata")
			return &userInfo, nil
		}
		log.Println("Redis failed,turn to DB query:", err)
	} else if err != redis.Nil {
		// 不是缓存未命中，而是 Redis 出错
		log.Println("Redis Get error:", err)
	}

	// Redis 无缓存或解析失败，查询数据库
	var user model.User
	if err := s.DB.Preload("Addresses").Where("id = ?", req.UserId).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "user is NOT FOUND!")
	}

	// 构建地址列表
	var addresses []*pb.Address
	for _, addr := range user.Addresses {
		addresses = append(addresses, &pb.Address{
			Id:       addr.ID,
			Label:    addr.Label,
			Province: addr.Province,
			City:     addr.City,
			District: addr.District,
			Detail:   addr.Detail,
		})
	}

	// 构建用户信息结构体
	userInfo := &pb.UserInfo{
		UserId:    fmt.Sprintf("%d", user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Addresses: addresses,
		IsVip:     user.IsVIP,
		VipLevel:  user.VIPLevel,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	// 序列化存入 Redis，设置过期时间
	data, err := json.Marshal(userInfo)
	if err == nil {
		err = config.RedisClient.Set(ct, cacheKey, data, time.Hour*24).Err()
		if err != nil {
			log.Println("Redis set error:", err)
		}
	} else {
		log.Println("user info xuliehua error:", err)
	}

	return userInfo, nil
}

// 实现更新地址簿的逻辑
func (s *UserService) UpdateAddressBook(ctx context.Context, req *pb.AddressBookUpdateRequest) (*pb.GenericResponse, error) {
	var user model.User
	//查找用户失败
	if err := s.DB.Where("user_id = ?", req.UserId).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "user is NOT FOUND!")
	}

	// 先删除旧地址
	if err := s.DB.Where("user_id = ?", user.ID).Delete(&model.Address{}).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "clear old address Failed!")
	}

	// 插入新地址
	for _, addr := range req.Addresses {
		newAddr := model.Address{
			ID:       addr.Id,
			Label:    addr.Label,
			Province: addr.Province,
			City:     addr.City,
			District: addr.District,
			Detail:   addr.Detail,
		}
		s.DB.Create(&newAddr)
	}

	return &pb.GenericResponse{
		Success: true,
		Message: "AddressBook update success!",
	}, nil
}
