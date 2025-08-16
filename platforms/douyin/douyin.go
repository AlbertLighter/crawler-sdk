package douyin

import (
	"crawler-sdk/pkg/http"
	"fmt"

	"resty.dev/v3"
)

// AuthDetails holds the temporary credentials for uploading.
// It's derived from the /web/api/media/upload/auth/v5/ endpoint response.

// DyClient 是抖音平台的客户端
type DyClient struct {
	client *resty.Client
}

// New 创建一个新的抖音客户端实例
// cookie: 用于身份验证的抖音 cookie 字符串
func New(cookie string) *DyClient {
	c := &DyClient{
		client: http.NewClient(cookie),
	}
	c.client.SetProxy("http://127.0.0.1:8888")
	c.client.AddRequestMiddleware(headers)
	return c
}

// GetVideos 实现了获取视频的逻辑
func (d *DyClient) GetVideos(query string) ([]string, error) {
	// fmt.Printf("Getting user [%s] from Douyin...\n", query)
	// 在这里实现具体的视频获取逻辑
	// 示例: d.client.R().Get(...)
	return []string{"douyin_video_1.mp4", "douyin_video_2.mp4"}, nil
}

func headers(c *resty.Client, req *resty.Request) error {
	fmt.Println("headers")
	// req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("sec-ch-ua-platform", "Windows")
	// req.Header.Set("sec-ch-ua", `"Not;A=Brand";v="99", "Microsoft Edge";v="139", "Chromium";v="139"`)
	// req.Header.Set("sec-ch-ua-mobile", "?0")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")
	// req.Header.Set("Accept", "application/json, text/plain, */*")
	// req.Header.Set("Origin", "https://www.xiaohongshu.com")
	// req.Header.Set("Sec-Fetch-Site", "same-site")
	// req.Header.Set("Sec-Fetch-Mode", "cors")
	// req.Header.Set("Sec-Fetch-Dest", "empty")
	return nil
}
