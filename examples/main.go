package main

import (
	"context"
	"crawler-sdk/platforms/xhs"
	"fmt"
	"log"
	"time"
)

func main() {
	// 平台特定的 cookies
	const xhsCookie = `abRequestId=a9e19a3c-ec62-5d58-82d1-11780180d594; a1=1989d69b653ej0uroayuyvsm4vdy26twup2h6vzan50000384309; webId=936fb47d51c175f4dec0400d43cb1af2; gid=yjYjfKjSDy7jyjYjfKjDKUhu2qdV87Fl0v7vJ0IhM64hfx28vJK9EK888qY4q8j8Y84JiYSf; x-user-id-creator.xiaohongshu.com=60de949b0000000001000c2f; customer-sso-sid=68c51753761449804526765463mdqkwv3gzu3e2d; customerClientId=913136795330800; access-token-creator.xiaohongshu.com=customer.creator.AT-68c5175376144980452676570jcbajjyvvfeitaq; galaxy_creator_session_id=P5r2edI91WGTWafcQqcfekMM05fSN8bK1sAb; galaxy.creator.beaker.session.id=1754987635895010773996; webBuild=4.75.3; web_session=040069b1aad7d3565c79251aa83a4b22b8638f; xsecappid=xhs-pc-web; loadts=1755165211163; acw_tc=0a50885317551652109482088e3f446d3b80b98308dd3884df9ecca40fd6d3; websectiga=3633fe24d49c7dd0eb923edc8205740f10fdb18b25d424d2a2322c6196d2a4ad; sec_poison_id=1003c47d-db6b-4c3d-bc0d-100c0b71ac3c; unread={%22ub%22:%22689d42df000000001b035010%22%2C%22ue%22:%22688aa2bf000000000302629f%22%2C%22uc%22:23}` // 替换为你的小红书 cookie
	// const douyinCookie = "your_douyin_cookie_here"                  // 替换为你的抖音 cookie

	// 示例用户ID
	const xhsUserID = "6763862e000000001500472d"

	// 1. 独立实例化每个平台的客户端
	xhsClient := xhs.New(xhsCookie)
	// douyinClient := douyin.New(douyinCookie)

	// 2. 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. 调用小红书的用户信息接口
	fmt.Printf("--- 正在获取小红书用户 [%s] 的信息 ---\n", xhsUserID)
	// userProfile, err := xhsClient.GetUserDetail(ctx, xhsUserID)
	// if err != nil {
	// 	log.Fatalf("获取小红书用户信息失败: %v", err)
	// }

	notes, err := xhsClient.GetUserNotes(ctx, xhsUserID, 0, 10)
	if err != nil {
		log.Fatalf("获取小红书用户笔记失败: %v", err)
	}
	fmt.Println(notes)

	// 4. 格式化为JSON并打印结果
	// prettyJSON, err := json.MarshalIndent(userProfile, "", "  ")
	// if err != nil {
	// 	log.Fatalf("格式化JSON失败: %v", err)
	// }

	// fmt.Println("成功获取到用户及笔记信息:")
	// fmt.Println(string(prettyJSON))

	// 5. 调用抖音的视频接口
	// fmt.Println("\n--- 抖音平台 (示例) ---")
	// videosDouyin, err := douyinClient.GetVideos("golang教程")
	// if err != nil {
	// 	log.Fatalf("从抖音获取视频失败: %v", err)
	// }
	// fmt.Printf("成功从抖音获取到视频: %v\n", videosDouyin)
}
