package main

import (
	"fmt"

	"github.com/spf13/viper"
	"turing.com/push/grpc/grpcserver"
	"turing.com/push/socket"
)

func main() {
	initConfig()

	// grpc
	go grpcserver.Init()
	// 初始化socket
	socket.InitSocketManager()
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))

}
