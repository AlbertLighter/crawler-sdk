package douyin

import (
	"fmt"

	"resty.dev/v3"
)

// Client 是抖音平台的客户端
type Client struct {
	client *resty.Client
}

// New 创建一个新的抖音客户端实例
func New(client *resty.Client) *Client {
	return &Client{client: client}
}

// GetVideos 实现了获取视频的逻辑
func (d *Client) GetVideos(query string) ([]string, error) {
	fmt.Printf("正在从抖音获取查询 [%s] 的视频...\n", query)
	// 在这里实现具体的视频获取逻辑
	return []string{"douyin_video_1.mp4", "douyin_video_2.mp4"}, nil
}