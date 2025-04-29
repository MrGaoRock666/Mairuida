package handler

import (
	"github.com/MrGaoRock666/Mairuida/order_service/model" // 数据模型
	pb "github.com/MrGaoRock666/Mairuida/order_service/pb" // proto 生成的 Go 包

	// 数据库配置
	"context" // 上下文
	"errors"
	"fmt"
	"log"
	"strconv" //字符串转整数
	"time"    // 时间处理

	"github.com/google/uuid" // 用于生成唯一订单号
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb" // 处理 protobuf 时间类型
	"gorm.io/gorm"

	"encoding/json"

	"github.com/MrGaoRock666/Mairuida/order_service/config" // 确保导入 config 包
	"github.com/go-redis/redis/v8"
)

// 定义OrderService结构体，实现pb.OrderService接口
type OrderService struct {
	DB                                 *gorm.DB               //注入数据库连接
	RedisClient                        *redis.Client          //注入redis连接
	ShippingConfig                     *config.ShippingConfig //运费计算配置
	pb.UnimplementedOrderServiceServer                        // 嵌套接口默认实现
}

// parseUint 将字符串安全转换为 uint 类型
func parseUint(s string) (uint, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	//如果 UserID 传了乱七八糟的字符串，会默默变成 0，需要检查
	if err != nil {
		return 0, err
	}
	return uint(u), nil
}

// 实现创建订单逻辑
func (s *OrderService) CreateOrder(ct context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {

	//使用事务来保证向数据库插入数据过程的原子性
	//开启事务
	tx := s.DB.Begin()

	//如果中间代码 panic 或其他异常，事务不会自动回滚，所以应该用 defer 加保护。
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 继续抛出
		}
	}()

	// 生成唯一订单号
	//UUID,即通用唯一标识符，是一种不依赖数据库自增 ID的、全局唯一的标识符，它的格式是一个 128 位的数字，通常以 32 个十六进制数字表示
	//在分布式系统里，每个服务实例自己生成 ID，如果还用数据库的自增 ID，多个实例可能会冲突
	//UUID 就是在分布式系统中生成“不会重复的订单号、用户 ID、请求 ID”的神器，安全、独立、无冲突
	//常用 Google 提供的 github.com/google/uuid 包，它默认生成的是 v4 类型
	//一旦使用了UUID，全世界任何一台服务器都不会生成一样的 ID，可靠且无需协调
	orderID := uuid.New().String()

	// 将字符串转成 uint 类型
	userID, err := parseUint(req.UserId)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("invalid user_id format")
	}
	// 订单模型实例
	order := &model.Order{
		OrderID:               orderID,
		UserID:                userID,
		SenderAddress:         req.SenderAddressId, // 注意：这里是地址 ID，后期可查地址详情
		ReceiverAddress:       req.ReceiverAddressId,
		ItemName:              req.ItemName,
		Weight:                req.Weight,
		Volume:                req.Volume,
		LogisticsCompany:      req.LogisticsCompany,
		PreferredDeliveryTime: req.PreferredDeliveryTime.AsTime(), // 转换时间格式
		IsUrgent:              req.IsUrgent,
		Status:                model.StatusCreated, // 初始状态为已创建
		CreatedAt:             time.Now(),          // 创建时间
		UpdatedAt:             time.Now(),          // 更新时间
	}

	// 保存订单到数据库，如果出错，事务是需要回滚的
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//结束事务
	tx.Commit()

	// 返回响应
	return &pb.OrderResponse{
		OrderId:   order.OrderID,
		Message:   "Create order Success",
		CreatedAt: timestamppb.New(order.CreatedAt), // 转为 protobuf 时间戳格式
	}, nil
}

// 实现查询订单详情逻辑
func (s *OrderService) GetOrder(ct context.Context, req *pb.OrderIDRequest) (*pb.OrderInfo, error) {
	// 先从 Redis 中查询订单信息
	orderJSON, err := s.RedisClient.Get(ct, "order:"+req.OrderId).Result()
	if err == nil {
		var orderInfo pb.OrderInfo
		// 将从 Redis 中获取的 JSON 字符串反序列化为 pb.OrderInfo 类型
		err := json.Unmarshal([]byte(orderJSON), &orderInfo)
		// 如果反序列化成功，没有错误
		if err == nil {
			return &orderInfo, nil
		}
	}

	var order model.Order

	// 从数据库查找订单--区分逻辑错误与系统错误
	if err := s.DB.First(&order, "order_id = ?", req.OrderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//找不到订单
			return nil, status.Errorf(codes.NotFound, "Order not found")
		}
		//系统错误(数据库错误)
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	// 构造返回消息
	orderInfo := &pb.OrderInfo{
		OrderId:               order.OrderID,
		UserId:                strconv.FormatUint(uint64(order.UserID), 10), // 将字符串转成 uint 类型
		SenderAddress:         order.SenderAddress,
		ReceiverAddress:       order.ReceiverAddress,
		ItemName:              order.ItemName,
		Weight:                order.Weight,
		Volume:                order.Volume,
		Status:                pb.OrderStatus(order.Status),
		CreatedAt:             timestamppb.New(order.CreatedAt),
		LogisticsCompany:      order.LogisticsCompany,
		PreferredDeliveryTime: timestamppb.New(order.PreferredDeliveryTime),
		IsUrgent:              order.IsUrgent,
		UpdatedAt:             timestamppb.New(order.UpdatedAt),
	}

	// 将订单信息存入 Redis
	orderJSONBytes, err := json.Marshal(orderInfo)
	if err == nil {
		// 将 []byte 转换为 string
		orderJSON := string(orderJSONBytes)
		s.RedisClient.Set(ct, "order:"+req.OrderId, orderJSON, 0)
	}

	return orderInfo, nil
}

// 模块核心：运费估算逻辑
// 经搜集资料，物流行业计算运费的大致行业标准如下
// 所以咱们也需要基于这点来设计运费估算算法
// 因为一个包裹可能存在(大而轻，小而重的现象)
// 行业通用标准是取"实际重量"和"体积重量"中的较大值作为计费重量
// 体积重量计算公式：长(cm)×宽(cm)×高(cm)÷体积系数（快递常用5000-6000，物流常用4800）--本项目暂定体积系数取5000
// 总费用 = 起步价 + (计费重量 × 单价/公里) × 距离
// 此处先实现一个初级算法，后续根据实际情况做改动
func (s *OrderService) EstimateCost(ct context.Context, req *pb.EstimateRequest) (*pb.EstimateResponse, error) {
	if s.ShippingConfig == nil {
		return nil, errors.New("shipping config is not initialized")
	}
	if s.RedisClient == nil {
		return nil, errors.New("Redis client is not initialized")
	}
	// 生成缓存键
	cacheKey := fmt.Sprintf("estimate:%f:%f:%f:%t:%t:%d", req.DistanceKm, req.Weight, req.Volume, req.IsUrgent, req.IsDelayed, req.TransportMethod)

	// 先从 Redis 中查询估算结果
	totalCost, err := s.RedisClient.Get(ct, cacheKey).Float64()
	if err != nil {
		log.Printf("Failed to get estimate cost from Redis: %v", err)
	}

	// 常量部分（可以后期做配置、动态调整），这里我先做个假设
	// 常量最好写成浮点数，避免浮点运算出现问题
	basePrice := s.ShippingConfig.BasePrice             // 起步价，单位：元
	pricePerKgPerKm := s.ShippingConfig.PricePerKgPerKm // 每公斤每公里费用，单位：元
	volumeRate := s.ShippingConfig.VolumeRate           //体积系数：本项目暂定为5000
	urgentSurcharge := s.ShippingConfig.UrgentSurcharge // 加急运输附加系数
	delayDiscount := s.ShippingConfig.DelayDiscount     // 延迟运输折扣系数

	// 获取运输方式对应的费率
	transportRate, ok := s.ShippingConfig.TransportRates[pb.TransportMethod_name[int32(req.TransportMethod)]]
	if !ok {
		return nil, errors.New("invalid transport method")
	}

	// 获取请求内带的重量和体积
	realWeight := req.Weight
	volume := req.Volume // 单位：立方米 m³

	// 计算体积重量（体积换算成重量）
	volumeWeight := volume * volumeRate

	// 计费重量取两者中的最大值
	billingWeight := realWeight
	if volumeWeight > realWeight {
		billingWeight = volumeWeight
	}

	// 校验逻辑，防止有人传0或者负数
	if req.DistanceKm <= 0 || req.Weight <= 0 {
		return nil, errors.New("invalid request parameters")
	}

	// 计算总费用：起步价 + (计费重量 × 单价/公里) × 距离
	totalCost = basePrice + (billingWeight * pricePerKgPerKm * req.DistanceKm)

	// 根据运输方式调整运费
	totalCost *= transportRate

	// 根据是否加急或延迟调整运费
	if req.IsUrgent {
		if req.TransportMethod == pb.TransportMethod_AIR {
			// 空运与加急绑定
			totalCost *= urgentSurcharge
		} else {
			totalCost *= urgentSurcharge
		}
	} else if req.IsDelayed {
		totalCost *= delayDiscount
		// 后续可以添加碳积分记录逻辑
	}

	// 将估算结果存入 Redis,利用set数据结构
	s.RedisClient.Set(ct, cacheKey, totalCost, 0)

	// 返回估算结果
	return &pb.EstimateResponse{
		EstimatedCost: totalCost,
	}, nil
}
