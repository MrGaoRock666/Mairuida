// 实现各RPC服务
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/MrGaoRock666/Mairuida/user_service/config"
	"github.com/MrGaoRock666/Mairuida/user_service/middleware"
	"github.com/MrGaoRock666/Mairuida/user_service/model"
	"github.com/MrGaoRock666/Mairuida/user_service/pb"
	"github.com/MrGaoRock666/Mairuida/user_service/util"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// 定义UserService结构体，实现pb.UserService接口
type UserService struct {
	DB                                *gorm.DB //注入数据库连接
	pb.UnimplementedUserServiceServer          //为没有实现的 RPC 方法提供默认的“未实现”报错处理，从而避免服务 crash
}

// 实现用户注册逻辑
// context 是 Go 语言中用于管理跨 API 边界和进程间请求范围数据、取消信号和超时的重要包
func (s *UserService) RegisterUser(ct context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//使用事务来保证向数据库插入数据过程的原子性
	//开启事务
	tx := s.DB.Begin()

	// 加密用户密码(bcrypt 加密逻辑)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Password generate error" + err.Error(),
		}, nil
	}

	//创建用户模型实例
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword), // 保存加密后的密码
		Email:    req.Email,
		Phone:    req.Phone,
	}
	//以下几条逻辑，都注定了，一个用户名，一个邮箱，一个电话号码，只可以注册一个用户。

	//查询用户名是否已经存在
	if err := tx.Where("username = ?", req.Username).First(&user).Error; err == nil {
		// 用户名已存在
		tx.Rollback()
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
	if err := tx.Create(&user).Error; err != nil {
		log.Println("Register failed:", err)
		//返回失败响应
		return &pb.RegisterResponse{
			Success: false,
			Message: "Register failed:" + err.Error(),
		}, nil
	}

	//结束事务
	tx.Commit()

	// 注册成功后，发送 Kafka 消息
	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to marshal user data: %v", err)
	} else {
		err = util.SendMessage(ct, string(userJson))
		if err != nil {
			log.Printf("Failed to send Kafka message: %v", err)
		}
	}

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

	//如果用户名错误，登录失败(这里只查用户名)
	if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "The user name is MOT exist")
	}

	// 使用 bcrypt 对比密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "密码错误")
	}

	// 模拟生成JWT token（TODO:后面再实现真正的JWT）
	//token := "fake-jwt-token-" + fmt.Sprintf("%d", user.ID)

	//登录成功，生成JWT Token
	token, err := util.GenerateJWT(uint64(user.ID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Token Failed!")
	}

	// 返回登录响应
	return &pb.LoginResponse{
		Token:  token,
		UserId: fmt.Sprintf("%d", user.ID),
	}, nil
}

// 实现获取用户信息逻辑（含 Redis 缓存）
func (s *UserService) GetUserInfo(ct context.Context, req *pb.UserIDRequest) (*pb.UserInfo, error) {
	// // 构建 Redis 的 key
	// cacheKey := fmt.Sprintf("user_info:%s", req.UserId)

	// //尝试从 Redis 获取缓存
	// cached, err := config.RedisClient.Get(ct, cacheKey).Result()
	// if err == nil {
	// 	// Redis 中有数据
	// 	var userInfo pb.UserInfo
	// 	if err := json.Unmarshal([]byte(cached), &userInfo); err == nil {
	// 		log.Println("GetUserInfo: return Redis cachedata")
	// 		return &userInfo, nil
	// 	}
	// 	log.Println("Redis failed,turn to DB query:", err)
	// } else if err != redis.Nil {
	// 	// 不是缓存未命中，而是 Redis 出错
	// 	log.Println("Redis Get error:", err)
	// }

	// // Redis 无缓存或解析失败，查询数据库
	// var user model.User

	// // 将 string 类型的 UserId 转为 uint
	// uid, err := strconv.ParseUint(req.UserId, 10, 64)

	// 有了JWT拦截器后，我无需通过查询redis和mysql来验证用户身份了。
	// Handler 代码 不再需要验证用户身份，只需从 context 中提取用户 ID 并进行后续处理。
	// 从 context 中获取 user_id

	var user model.User

	uidStr, ok := middleware.GetUserIDFromContext(ct)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user_id not found in context")
	}

	cacheKey := fmt.Sprintf("user_info:%s", uidStr)

	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	//注意：这里对软删除行为有一个布尔控制，如果之前注销过账户，布尔标记为true，则不会跑数据库的查询逻辑
	if err := s.DB.Preload("Addresses").Where("id = ? AND is_deleted = ?", uid, false).First(&user).Error; err != nil {
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
func (s *UserService) UpdateAddressBook(ct context.Context, req *pb.AddressBookUpdateRequest) (*pb.GenericResponse, error) {
	//var user model.User
	// //查找用户失败
	// if err := s.DB.Where("id = ? AND is_deleted = ?", req.UserId, false).First(&user).Error; err != nil {
	// 	return nil, status.Errorf(codes.NotFound, "user is NOT FOUND!")
	// }

	// 有了JWT拦截器后，从 context 获取当前用户的 user_id
	userID, ok := middleware.GetUserIDFromContext(ct)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context")
	}

	var user model.User
	// 查找用户失败
	if err := s.DB.Where("id = ? AND is_deleted = ?", userID, false).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found!")
	}

	// 开启事务
	tx := s.DB.Begin()

	// 先删除旧地址
	if err := tx.Where("user_id = ?", user.ID).Delete(&model.Address{}).Error; err != nil {
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
			UserID:   user.ID,
		}
		//如果失败，需要进行事务回滚
		if err := tx.Create(&newAddr).Error; err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "insert address Failed!")
		}
	}

	// 提交事务
	tx.Commit()

	//返回一个通用响应
	return &pb.GenericResponse{
		Success: true,
		Message: "AddressBook update success!",
	}, nil
}

// 生成一个 6 位随机验证码（字符串形式）
func generate6DigitCode() string {
	rand.Seed(time.Now().UnixNano())               // 使用当前时间作为随机数种子
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // 格式化为 6 位数字，不足前面补 0
}

// 实现生成并发送验证码的逻辑
func (s *UserService) SendLoginCode(ct context.Context, req *pb.LoginCodeRequest) (*pb.GenericResponse, error) {
	code := generate6DigitCode()          // 调用自定义函数生成 6 位随机验证码
	cacheKey := "login_code:" + req.Phone // Redis 中存储验证码的 key，便于按手机号查找

	// 将验证码存入 Redis(因为验证码是“临时性的认证凭证”)，设置过期时间为 5 分钟
	// 用 Redis.Set(key(电话号码), value(短信验证码), 5min) 存一条数据
	err := config.RedisClient.Set(ct, cacheKey, code, 5*time.Minute).Err()
	if err != nil {
		// 存储失败，返回内部错误
		return nil, status.Errorf(codes.Internal, "send code failed")
	}

	// TODO注释：明确告诉开发者这里还没写完，要补充
	// TODO: 这里应接入第三方短信平台（如阿里云 / 腾讯云）实现验证码发送
	// 目前只是模拟成功返回
	return &pb.GenericResponse{Success: true, Message: "code sent"}, nil
}

// 实现使用验证码登录的逻辑
func (s *UserService) LoginByCode(ctx context.Context, req *pb.LoginByCodeRequest) (*pb.LoginResponse, error) {
	cacheKey := "login_code:" + req.Phone // 从 Redis 获取验证码的 key
	//后端用手机号构造 Redis key 再 .Get() 查缓存，看看用户输的和 Redis 里存的是否一致。如果一致，就登录成功；否则失败。
	code, err := config.RedisClient.Get(ctx, cacheKey).Result()

	// 验证码错误或 Redis 查询失败
	if err != nil || code != req.Code {
		return nil, status.Errorf(codes.Unauthenticated, "invalid code")
	}

	var user model.User
	// 用手机号从数据库查找用户
	if err := s.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "User NOT FOUND")
	}

	// 模拟生成 token（TODO:后续可接入 JWT）
	token := "fake-token-" + fmt.Sprintf("%d", user.ID)

	// 返回登录响应
	return &pb.LoginResponse{Token: token, UserId: fmt.Sprintf("%d", user.ID)}, nil
}

// 实现注销用户逻辑
func (s *UserService) DeleteUserAccount(ct context.Context, req *pb.UserIDRequest) (*pb.GenericResponse, error) {
	// 从 context 获取当前用户的 user_id
	userID, ok := middleware.GetUserIDFromContext(ct)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context")
	}

	var user model.User

	//根据请求的用户ID查找用户，使用 First 查询，若没有找到会返回 gorm.ErrRecordNotFound 错误
	if err := s.DB.First(&user, userID).Error; err != nil {
		// 如果没有找到用户，返回错误提示
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.GenericResponse{
				Success: false,
				Message: "The user is NOT EXIST", // 用户不存在时返回错误信息
			}, nil
		}
		// 发生其他错误时，直接返回错误
		return nil, err
	}

	//判断用户是否已经注销
	//为什么要有这一步？因为我们这个系统设计的是软删除机制
	//也就是说数据库依然有备查数据的
	if user.IsDeleted {
		// 如果用户已经注销，直接返回提示信息
		return &pb.GenericResponse{
			Success: false,
			Message: "The user has already deleted!", // 如果用户已注销，提示用户已注销
		}, nil
	}

	// 开启事务，保证程序的原子性，鲁棒性和可维护性
	tx := s.DB.Begin()

	//进行软删除，标记用户为已注销
	user.IsDeleted = true
	// 更新数据库中的用户记录，将 IsDeleted 字段更新为 true
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err // 如果保存失败，返回错误
	}

	// TODO: 如果未来删除更多资源，这里都可以继续 tx.xxx()

	// 提交事务
	tx.Commit()

	// 删除 Redis 中的登录状态（如果有的话）
	// 通过 Redis key 来删除对应的用户登录信息
	// redis非事务处理，但记录日志
	redisKey := fmt.Sprintf("login:user:%d", user.ID) // 生成 Redis key
	if err := config.RedisClient.Del(ct, redisKey).Err(); err != nil {
		// 如果删除 Redis 登录信息失败，打印警告日志
		log.Printf("Redis delete error:%v", err)
	}

	// 返回通用响应，成功注销
	return &pb.GenericResponse{
		Success: true,
		Message: "deleted successfully!thank for your use", // 用户成功注销时返回成功消息
	}, nil
}

// 实现恢复已注销用户的逻辑
// 恢复用户逻辑
func (s *UserService) RestoreUserAccount(ctx context.Context, req *pb.UserRestoreRequest) (*pb.GenericResponse, error) {
	// 从 context 获取当前用户的 user_id
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context")
	}

	var user model.User

	// 查询用户是否存在（无论是否被注销）
	if err := s.DB.First(&user, userID).Error; err != nil {
		return &pb.GenericResponse{
			Success: false,
			Message: "User NOT FOUND!",
		}, nil
	}

	// 用户未注销，无需恢复
	if !user.IsDeleted {
		return &pb.GenericResponse{
			Success: false,
			Message: "User is not deleted!Can NOT restore",
		}, nil
	}

	// 开启事务
	tx := s.DB.Begin()

	// 恢复账号
	user.IsDeleted = false
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Restore failed: %v", err)
	}

	tx.Commit()

	// 清理 Redis 缓存（避免旧状态）
	cacheKey := fmt.Sprintf("user_info:%d", user.ID)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Redis delete failed: %v", err)
	}

	//返回恢复成功响应
	return &pb.GenericResponse{
		Success: true,
		Message: "Account restored successfully!",
	}, nil
}
