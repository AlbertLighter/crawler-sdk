package douyin

import (
	"crawler-sdk/pkg/http"

	"resty.dev/v3"
)

// DyClient 是抖音平台的客户端
type DyClient struct {
	client *resty.Client
}

// New 创建一个新的抖音客户端实例
// cookie: 用于身份验证的抖音 cookie 字符串
func New(cookie string) *DyClient {
	return &DyClient{
		client: http.NewClient(cookie),
	}
}

// GetVideos 实现了获取视频的逻辑
func (d *DyClient) GetVideos(query string) ([]string, error) {
	// fmt.Printf("Getting user [%s] from Douyin...\n", query)
	// 在这里实现具体的视频获取逻辑
	// 示例: d.client.R().Get(...)
	return []string{"douyin_video_1.mp4", "douyin_video_2.mp4"}, nil
}
