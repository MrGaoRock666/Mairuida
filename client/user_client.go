package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "Mairuida/user_service/pb"
)

func main() {
	// 连接到 gRPC 服务
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx := context.Background()

	// 1. 测试注册
	regResp, err := client.RegisterUser(ctx, &pb.RegisterRequest{
		Username: "gbl",
		Password: "200510",
		Email:    "2801787642@qq.com",
		Phone:    "13919472452",
	})
	if err != nil {
		log.Println("注册失败:", err)
	} else {
		fmt.Println("注册响应:", regResp)
		// 获取注册成功后的 user_id，后续会用到这个 ID
		userId := regResp.UserId // 注册时返回的 user_id（例如 "2"）

		// 2. 测试登录
		loginResp, err := client.LoginUser(ctx, &pb.LoginRequest{
			Username: "gbl",
			Password: "200510",
		})
		if err != nil {
			log.Println("登录失败:", err)
		} else {
			// 登录成功，打印登录信息
			fmt.Println("登录成功,UserID:", loginResp.UserId, "Token:", loginResp.Token)

			// 3. 测试获取用户信息
			infoResp, err := client.GetUserInfo(ctx, &pb.UserIDRequest{
				UserId: userId, // 使用注册时返回的 user_id
			})
			if err != nil {
				log.Println("获取用户信息失败:", err)
			} else {
				fmt.Println("用户信息:", infoResp)
			}

			// 4. 测试更新地址簿
			_, err = client.UpdateAddressBook(ctx, &pb.AddressBookUpdateRequest{
				UserId: userId, // 使用注册时返回的 user_id
				Addresses: []*pb.Address{
					{
						Id:       "1", // 注意这里是字符串类型
						Label:    "家",
						Province: "甘肃省",
						City:     "兰州市",
						District: "安宁区",
						Detail:   "某小区1栋101",
					},
				},
			})
			if err != nil {
				log.Println("更新地址簿失败:", err)
			} else {
				fmt.Println("地址簿更新成功")
			}

			// 5. 模拟发送验证码（测试阶段直接返回成功）
			codeResp, err := client.SendLoginCode(ctx, &pb.LoginCodeRequest{
				Phone: "13888888888",
			})
			if err != nil {
				log.Println("发送验证码失败:", err)
			} else {
				fmt.Println("验证码发送响应:", codeResp)
			}

			// 6. 测试验证码登录（验证码为 Redis 中存的值，测试可用 fake）
			codeLoginResp, err := client.LoginByCode(ctx, &pb.LoginByCodeRequest{
				Phone: "13919472452",
				Code:  "685135", // 真实应从 Redis 中获取或输出
			})
			if err != nil {
				log.Println("验证码登录失败:", err)
			} else {
				fmt.Println("验证码登录成功,UserID:", codeLoginResp.UserId)
			}

			// 7. 测试注销账户
			//delResp, err := client.DeleteUserAccount(ctx, &pb.UserIDRequest{
				//UserId: userId, // 使用注册时返回的 user_id
			//})
			//if err != nil {
				//log.Println("注销失败:", err)
			//} else {
				//fmt.Println("注销响应:", delResp)
			//}
		}
	}
}
