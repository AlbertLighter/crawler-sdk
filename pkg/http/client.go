package http

import "resty.dev/v3"

// NewClient 创建并配置一个默认的resty客户端。
// cookie: 用于身份验证的cookie字符串
func NewClient(cookie string) *resty.Client {
	client := resty.New()

	// 在这里可以设置通用的配置，例如：
	// - User-Agent
	// - 超时时间
	// - 代理
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	if cookie != "" {
		client.SetHeader("Cookie", cookie)
	}
	client.SetProxy("http://127.0.0.1:8888")
	return client
}
