package main

import (
	"context"
	"crawler-sdk"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func main() {
	// 这是一个示例 cookie，请替换为你自己的有效 cookie
	const userCookie = "a1=188axxxxxxxx; web_session=04006bxxxxxxxx"
	// 这是一个示例用户ID
	const userID = "5b459a74e8ac2b5da3336423"

	// 1. 实例化主SDK客户端, 并传入cookie
	sdkClient := crawlersdk.NewClient(userCookie)

	// 2. 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. 调用小红书的用户信息接口
	fmt.Printf("--- 正在获取小红书用户 [%s] 的信息 ---\n", userID)
	userProfile, err := sdkClient.XHS.GetUser(ctx, userID)
	if err != nil {
		log.Fatalf("获取小红书用户信息失败: %v", err)
	}

	// 4. 格式化为JSON并打印结果
	prettyJSON, err := json.MarshalIndent(userProfile, "", "  ")
	if err != nil {
		log.Fatalf("格式化JSON失败: %v", err)
	}

	fmt.Println("成功获取到用户及笔记信息:")
	fmt.Println(string(prettyJSON))

	// 保留原始的GetVideos调用作为另一个示例
	fmt.Println("\n--- 抖音平台 (示例) ---")
	videosDouyin, err := sdkClient.Douyin.GetVideos("golang教程")
	if err != nil {
		log.Fatalf("从抖音获取视频失败: %v", err)
	}
	fmt.Printf("成功从抖音获取到视频: %v\n", videosDouyin)
}