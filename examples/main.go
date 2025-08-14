package main

import (
	"context"
	"crawler-sdk/platforms/douyin"
	"crawler-sdk/platforms/xhs"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func main() {
	// 平台特定的 cookies
	const xhsCookie = "a1=188axxxxxxxx; web_session=04006bxxxxxxxx" // 替换为你的小红书 cookie
	const douyinCookie = "your_douyin_cookie_here"                 // 替换为你的抖音 cookie

	// 示例用户ID
	const xhsUserID = "5b459a74e8ac2b5da3336423"

	// 1. 独立实例化每个平台的客户端
	xhsClient := xhs.New(xhsCookie)
	douyinClient := douyin.New(douyinCookie)

	// 2. 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. 调用小红书的用户信息接口
	fmt.Printf("--- 正在获取小红书用户 [%s] 的信息 ---\n", xhsUserID)
	userProfile, err := xhsClient.GetUser(ctx, xhsUserID)
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

	// 5. 调用抖音的视频接口
	fmt.Println("\n--- 抖音平台 (示例) ---")
	videosDouyin, err := douyinClient.GetVideos("golang教程")
	if err != nil {
		log.Fatalf("从抖音获取视频失败: %v", err)
	}
	fmt.Printf("成功从抖音获取到视频: %v\n", videosDouyin)
}
