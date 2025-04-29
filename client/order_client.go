// 订单服务测试客户端
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/MrGaoRock666/Mairuida/order_service/pb"
)

func main() {
	// 连接到订单服务
	conn, err := grpc.Dial("localhost:5002", grpc.WithInsecure()) // 假设订单服务在5002端口
	if err != nil {
		log.Fatalf("连接订单服务失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	ctx := context.Background()

	// 测试运费估算功能
	testEstimateCost(client, ctx)

	// 测试创建订单功能
	orderID := testCreateOrder(client, ctx)

	// 测试查询订单功能
	if orderID != "" {
		testGetOrder(client, ctx, orderID)
	}
}

// 测试运费估算
func testEstimateCost(client pb.OrderServiceClient, ctx context.Context) {
	fmt.Println("\n=== 开始测试运费估算 ===")

	resp, err := client.EstimateCost(ctx, &pb.EstimateRequest{
		DistanceKm:      120.5,                   // 距离120.5公里
		Weight:          3.2,                     // 重量3.2kg
		Volume:          0.0045,                  // 体积0.0045立方米(4.5升)
		TransportMethod: pb.TransportMethod_ROAD, // 指定运输方式为公路运输
	})
	if err != nil {
		log.Printf("运费估算失败: %v", err)
		return
	}

	fmt.Printf("运费估算结果:\n"+
		"距离: 120.5公里\n"+
		"重量: 3.2kg\n"+
		"体积: 0.0045m³\n"+
		"运输方式: 公路运输\n"+
		"估算费用: ￥%.2f元\n",
		resp.EstimatedCost)
}

// 测试创建订单
func testCreateOrder(client pb.OrderServiceClient, ctx context.Context) string {
	fmt.Println("\n=== 开始测试创建订单 ===")

	// 设置期望送达时间（当前时间+48小时）
	deliveryTime := time.Now().Add(48 * time.Hour)

	resp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId:                "123", // 测试用户ID
		SenderAddressId:       "addr_1",
		ReceiverAddressId:     "addr_2",
		ItemName:              "笔记本电脑",
		Weight:                3.2,
		Volume:                0.0045,
		LogisticsCompany:      "顺丰速运",
		PreferredDeliveryTime: timestamppb.New(deliveryTime),
		IsUrgent:              true,
	})
	if err != nil {
		log.Printf("创建订单失败: %v", err)
		return ""
	}

	fmt.Printf("订单创建成功:\n"+
		"订单ID: %s\n"+
		"消息: %s\n"+
		"创建时间: %v\n",
		resp.OrderId, resp.Message, resp.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))

	return resp.OrderId
}

// 测试查询订单
func testGetOrder(client pb.OrderServiceClient, ctx context.Context, orderID string) {
	fmt.Println("\n=== 开始测试查询订单 ===")

	resp, err := client.GetOrder(ctx, &pb.OrderIDRequest{
		OrderId: orderID,
	})
	if err != nil {
		log.Printf("查询订单失败: %v", err)
		return
	}

	fmt.Printf("订单查询结果:\n"+
		"订单ID: %s\n"+
		"用户ID: %s\n"+
		"物品名称: %s\n"+
		"重量: %.2fkg\n"+
		"体积: %.4fm³\n"+
		"物流公司: %s\n"+
		"状态: %s\n"+
		"是否加急: %t\n"+
		"创建时间: %v\n"+
		"期望送达时间: %v\n",
		resp.OrderId,
		resp.UserId,
		resp.ItemName,
		resp.Weight,
		resp.Volume,
		resp.LogisticsCompany,
		getStatusName(resp.Status),
		resp.IsUrgent,
		resp.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
		resp.PreferredDeliveryTime.AsTime().Format("2006-01-02 15:04:05"))
}

// 辅助函数：获取状态名称
func getStatusName(status pb.OrderStatus) string {
	switch status {
	case pb.OrderStatus_CREATED:
		return "已创建"
	case pb.OrderStatus_DISPATCHING:
		return "待配送"
	case pb.OrderStatus_IN_TRANSIT:
		return "运输中"
	case pb.OrderStatus_DELIVERED:
		return "已送达"
	case pb.OrderStatus_CANCELLED:
		return "已取消"
	default:
		return "未知状态"
	}
}
