package xhs

import (
	"fmt"

	"resty.dev/v3"
)

// Client 是小红书平台的客户端
type Client struct {
	client *resty.Client
}

// New 创建一个新的小红书客户端实例
func New(client *resty.Client) *Client {
	return &Client{client: client}
}

// GetVideos 实现了获取视频的逻辑
func (x *Client) GetVideos(query string) ([]string, error) {
	fmt.Printf("正在从小红书获取查询 [%s] 的视频...\n", query)
	// 在这里实现具体的视频获取逻辑
	return []string{"xhs_note_1", "xhs_note_2"}, nil
}