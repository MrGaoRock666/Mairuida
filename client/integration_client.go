package main

import (
	"fmt"
	"log"
	"time"

	orderConfig "github.com/MrGaoRock666/Mairuida/order_service/config"
	orderModel "github.com/MrGaoRock666/Mairuida/order_service/model"
	userModel "github.com/MrGaoRock666/Mairuida/user_service/model"
	"github.com/google/uuid"
)

func main() {
	orderDB := orderConfig.InitOrderDB()
	userDB := orderConfig.InitUserDB()

	// 插入测试用户数据
	testUser := userModel.User{
		Username: "hhl",
		Password: "78654399",
		Email:    "9996667788.com",
		Phone:    "1865492377",
	}
	result := userDB.Create(&testUser)
	if result.Error != nil {
		log.Fatalf("Failed to insert user: %v", result.Error)
	}

	// 插入测试地址数据
	testSenderAddress := userModel.Address{
		ID:       "9",
		UserID:   testUser.ID,
		Label:    "家",
		Province: "福建省",
		City:     "福州市",
		District: "闽侯县",
		Detail:   "福建师范大学",
	}
	result = userDB.Create(&testSenderAddress)
	if result.Error != nil {
		log.Fatalf("Failed to insert sender address: %v", result.Error)
	}

	testReceiverAddress := userModel.Address{
		ID:       "10",
		UserID:   testUser.ID,
		Province: "安徽省",
		City:     "黄山市",
		District: "黄山市区",
		Detail:   "汤口村",
	}
	result = userDB.Create(&testReceiverAddress)
	if result.Error != nil {
		log.Fatalf("Failed to insert receiver address: %v", result.Error)
	}

	// 生成唯一的 order_id
	uniqueOrderID := uuid.New().String()

	// 插入测试订单数据
	testOrder := orderModel.Order{
		OrderID:               uniqueOrderID, // 设置唯一的 order_id
		SenderAddress:         "9",
		ReceiverAddress:       "10",
		PreferredDeliveryTime: time.Now(),
	}
	result = orderDB.Create(&testOrder)
	if result.Error != nil {
		log.Fatalf("Failed to insert order: %v", result.Error)
	}

	// 获取订单信息
	var order orderModel.Order
	result = orderDB.First(&order)
	if result.Error != nil {
		log.Fatalf("Failed to get order: %v", result.Error)
	}

	// 获取发件地址
	senderAddress, err := orderConfig.GetAddressByID(userDB, order.SenderAddress)
	if err != nil {
		fmt.Printf("Failed to get sender address: %v\n", err)
	} else {
		fmt.Printf("Sender Address: %+v\n", senderAddress)
	}

	// 获取收件地址
	receiverAddress, err := orderConfig.GetAddressByID(userDB, order.ReceiverAddress)
	if err != nil {
		fmt.Printf("Failed to get receiver address: %v\n", err)
	} else {
		fmt.Printf("Receiver Address: %+v\n", receiverAddress)
	}
}
