package main

import (
	"fmt"
	"yikou-ai-go-teach/wire"
)

func main() {
	// 初始化 Web 服务器
	h, err := wire.InitializeApp()
	if err != nil {
		panic(fmt.Errorf("注入失败：: %w", err))
	}
	// 启动服务器
	h.Spin()
}
