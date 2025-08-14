package crawlersdk

import (
	"crawler-sdk/pkg/http"
	"crawler-sdk/platforms/douyin"
	"crawler-sdk/platforms/xhs"
)

// Client是SDK的主客户端，封装了对所有平台的操作
type Client struct {
	Douyin *douyin.Client
	XHS    *xhs.Client
}

// NewClient 创建一个新的SDK客户端实例
// cookie: 用于身份验证的cookie字符串
func NewClient(cookie string) *Client {
	httpClient := http.NewClient(cookie)

	return &Client{
		Douyin: douyin.New(httpClient),
		XHS:    xhs.New(httpClient),
	}
}
