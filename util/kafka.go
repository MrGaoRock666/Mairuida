package util

import (
	"context"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

// InitKafka 初始化 Kafka 写入器，用于向指定的 Kafka 主题发送消息。
// 参数 brokers 是一个字符串切片，包含 Kafka 集群的地址列表，例如 ["localhost:9092"]。
func InitKafka(brokers []string) {
	// 创建一个新的 kafka.Writer 实例，并将其赋值给全局变量 kafkaWriter。
	// 该写入器用于向 Kafka 主题发送消息。
	kafkaWriter = &kafka.Writer{
		// Addr 指定 Kafka 集群的地址，使用 kafka.TCP 函数将 brokers 切片转换为网络地址。
		Addr: kafka.TCP(brokers...),
		// Topic 指定要发送消息的 Kafka 主题，这里设置为 "user_register_topic"，
		// 表示该写入器将消息发送到用户注册通知主题。
		Topic: "user_register_topic",
	}
}

// SendMessage 发送消息到 Kafka 主题。
// 参数 ctx 为上下文，用于控制消息发送操作的生命周期，可设置超时、取消等。
// 参数 msg 为要发送的消息内容，类型为字符串。
// 返回值为 error 类型，若消息发送过程中出现错误，则返回相应的错误信息；若发送成功，则返回 nil。
func SendMessage(ctx context.Context, msg string) error {
	// 检查 kafkaWriter 是否为 nil，若为 nil 则表示 Kafka 写入器未初始化，
	// 此时不进行消息发送操作，直接返回 nil。
	if kafkaWriter == nil {
		return nil
	}
	// 调用 kafkaWriter 的 WriteMessages 方法将消息发送到 Kafka 主题。
	// 将传入的字符串消息转换为字节切片作为消息的 Value 部分。
	return kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: []byte(msg),
	})
}

// CloseKafka 关闭 Kafka 写入器
func CloseKafka() error {
	if kafkaWriter == nil {
		return nil
	}
	return kafkaWriter.Close()
}
