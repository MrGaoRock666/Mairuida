package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// 创建一个新的 Kafka 读取器
	r := kafka.NewReader(kafka.ReaderConfig{
		// Brokers 指定 Kafka 集群的地址列表，这里使用本地默认端口 9092 的 Kafka 实例。
		Brokers: []string{"localhost:9092"},
		// Topic 指定要读取消息的 Kafka 主题，这里是用户注册通知主题。
		Topic: "user_register_topic",
		// GroupID 指定消费者所属的消费者组，用于实现消息的分区消费和负载均衡。
		GroupID: "user-register-group",
		// MinBytes 指定每次读取消息时，Kafka 服务器返回的最小字节数，这里设置为 10KB。
		MinBytes: 10e3, // 10KB
		// MaxBytes 指定每次读取消息时，Kafka 服务器返回的最大字节数，这里设置为 10MB。
		MaxBytes: 10e6, // 10MB
	})
	// 确保在 main 函数结束时关闭 Kafka 读取器，释放资源。
	defer r.Close()

	// 打印启动信息，提示 Kafka 消费者已成功启动。
	fmt.Println("Kafka consumer started...")

	// 使用无限循环持续从 Kafka 主题中读取消息。
	for {
		// 从 Kafka 主题中读取一条消息，使用 context.Background() 作为上下文。
		msg, err := r.ReadMessage(context.Background())
		// 检查读取消息过程中是否发生错误。
		if err != nil {
			// 若发生错误，使用 log.Fatalf 打印错误信息并终止程序。
			log.Fatalf("Failed to read message: %v", err)
		}
		// 模拟发送欢迎邮件的操作，打印接收到的消息内容。
		log.Printf("Received message: %s, sending welcome email...", string(msg.Value))
		// 此处可添加实际的发送邮件逻辑，例如调用邮件服务 API 发送欢迎邮件给新注册用户。
		// 示例代码：
		// sendWelcomeEmail(string(msg.Value))
	}
}
