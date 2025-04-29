package config

import (
	"os"

	"github.com/spf13/viper"
)

// ShippingConfig 定义运费计算相关的配置结构体
type ShippingConfig struct {
	BasePrice       float64 `mapstructure:"basePrice"`
	PricePerKgPerKm float64 `mapstructure:"pricePerKgPerKm"`
	VolumeRate      float64 `mapstructure:"volumeRate"`
	UrgentSurcharge float64 `mapstructure:"urgentSurcharge"`
	DelayDiscount   float64 `mapstructure:"delayDiscount"`
	//用map存储运输方式
	TransportRates map[string]float64 `mapstructure:"transportRates"`
}

// LoadShippingConfig 加载运费计算相关的配置信息。
// 该函数使用 viper 库从 JSON 配置文件中读取配置，并将 "shipping" 部分的配置
// 反序列化到 ShippingConfig 结构体中。
// 返回值：
//   - *ShippingConfig：包含运费计算配置的结构体指针。
//   - error：如果读取配置文件或反序列化过程中出现错误，返回相应的错误信息；否则返回 nil。
func LoadShippingConfig() (*ShippingConfig, error) {
	// 获取项目根目录的绝对路径
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 拼接配置文件的绝对路径，使用 = 赋值
	configPath := projectRoot + "/../config"

	// 设置要读取的配置文件的名称，不包含扩展名
	viper.SetConfigName("config")
	// 设置配置文件的类型为 JSON
	viper.SetConfigType("json")
	// 添加配置文件所在的目录，使用绝对路径
	viper.AddConfigPath(configPath)

	// 读取配置文件
	err = viper.ReadInConfig()
	// 如果读取配置文件时出现错误，直接返回 nil 和错误信息
	if err != nil {
		return nil, err
	}

	// 声明一个 ShippingConfig 类型的变量，用于存储反序列化后的配置信息
	var shippingConfig ShippingConfig
	// 将配置文件中 "shipping" 部分的配置信息反序列化到 shippingConfig 变量中
	err = viper.UnmarshalKey("shipping", &shippingConfig)
	// 如果反序列化过程中出现错误，直接返回 nil 和错误信息
	if err != nil {
		return nil, err
	}

	// 返回包含配置信息的 ShippingConfig 结构体指针和 nil 错误信息
	return &shippingConfig, nil
}
