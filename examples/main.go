package main

import (
	"crawler-sdk"
	"fmt"
	"log"
)

func main() {
	// 假设这是从外部传入的cookie
	// 在实际应用中, 这个cookie应该通过安全的方式获取和管理
	const userCookie = "your_actual_cookie_string_here"

	// 1. 实例化主SDK客户端, 并传入cookie
	sdkClient := crawlersdk.NewClient(userCookie)

	// 2. 通过主客户端调用各个平台的接口
	fmt.Println("--- 抖音平台 ---")

	videosDouyin, err := sdkClient.Douyin.GetVideos("golang教程")
	if err != nil {
		log.Fatalf("从抖音获取视频失败: %v", err)
	}
	fmt.Printf("成功从抖音获取到视频: %v\n", videosDouyin)

	fmt.Println("\n--- 小红书平台 ---")

	videosXHS, err := sdkClient.XHS.GetVideos("go语言学习")
	if err != nil {
		log.Fatalf("从小红书获取视频失败: %v", err)
	}
	fmt.Printf("成功从小红书获取到视频: %v\n", videosXHS)
}
