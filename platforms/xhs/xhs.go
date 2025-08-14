package xhs

import (
	"crawler-sdk/internal/crypto/xhs"
	"crawler-sdk/pkg/http"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"resty.dev/v3"
)

// Client 是小红书平台的客户端
type Client struct {
	client *resty.Client
}

// New 创建一个新的小红书客户端实例
// cookie: 用于身份验证的小红书 cookie 字符串
func New(cookie string) *Client {
	c := &Client{
		client: http.NewClient(cookie),
	}
	// c.client.AddRequestMiddleware(SignXYS)
	c.client.AddRequestMiddleware(headers)
	c.client.AddRequestMiddleware(SignXS)
	c.client.AddRequestMiddleware(SignXSC)
	c.client.AddRequestMiddleware(SignTraceID)
	c.client.SetProxy("http://127.0.0.1:8888")
	return c
}

// GET https://edith.xiaohongshu.com/api/sns/web/v1/user_posted?num=30&cursor=684ebc77000000002100a2c4&user_id=5d5c36ab0000000001008656&image_formats=jpg,webp,avif&xsec_token=ABn830XCOxiqnEyuW9NzS0hmuNr9Se3HJ3v-pZFRItuHo%3D&xsec_source=pc_feed HTTP/1.1
// Host: edith.xiaohongshu.com
// Connection: keep-alive
// sec-ch-ua-platform: "Windows"
// sec-ch-ua:
// sec-ch-ua-mobile: ?0
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0
// Accept: application/json, text/plain, */*
// Origin: https://www.xiaohongshu.com
// Sec-Fetch-Site: same-site
// Sec-Fetch-Mode: cors
// Sec-Fetch-Dest: empty
// Referer: https://www.xiaohongshu.com/

// GET https://edith.xiaohongshu.com/api/sns/web/v1/user_posted?cursor=&image_formats=jpg&image_formats=webp&image_formats=avif&num=30&user_id=6763862e000000001500472d HTTP/1.1
// Host: edith.xiaohongshu.com
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36
// Accept-Encoding: gzip, deflate

func headers(c *resty.Client, req *resty.Request) error {
	fmt.Println("headers")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-ch-ua", `"Not;A=Brand";v="99", "Microsoft Edge";v="139", "Chromium";v="139"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://www.xiaohongshu.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	return nil
}

func SignXYS(c *resty.Client, req *resty.Request) error {
	fmt.Println("SignXYS")
	xsc := xhs.XYS(req.URL, req.Header.Get("Cookie"))
	req.Header.Set("X-s", xsc)
	return nil
}

func SignXS(c *resty.Client, req *resty.Request) error {
	fmt.Println("SignXS")
	encryptor := xhs.NewXsEncrypt()
	cookie := c.Header().Get("Cookie")
	a1 := getCookieValue(cookie, "a1")
	ts := fmt.Sprintf("%d", time.Now().Unix()*1000)
	signedURL, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("URL解析失败: %w", err)
	}
	q := map[string]string{}
	for k, v := range req.QueryParams {
		q[k] = v[0]
	}
	qs, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %w", err)
	}
	signedURL.RawQuery = string(qs)
	fmt.Println(signedURL.String())
	xs, err := encryptor.EncryptXs(signedURL.String(), a1, ts, "xhs-pc-web")
	if err != nil {
		return fmt.Errorf("加密 'x-s' 失败: %w", err)
	}
	req.Header.Set("X-s", xs)
	req.Header.Set("X-t", ts)
	return nil
}

func SignXSC(c *resty.Client, req *resty.Request) error {
	fmt.Println("SignXSC")
	encryptor := xhs.XscEncrypt{}
	cookie := req.Header.Get("Cookie")
	a1 := getCookieValue(cookie, "a1")
	x1 := getCookieValue(cookie, "x1")
	x4 := getCookieValue(cookie, "x4")
	b1 := getCookieValue(cookie, "b1")
	xsc, err := encryptor.EncryptXsc(req.Header.Get("x-s"), req.Header.Get("X-t"), "xhs-pc-web", a1, x1, x4, b1)
	if err != nil {
		return fmt.Errorf("加密 'x-sc' 失败: %w", err)
	}
	xscStr := base64.StdEncoding.EncodeToString(xsc)
	req.Header.Set("X-S-Common", xscStr)
	return nil
}

func SignTraceID(c *resty.Client, req *resty.Request) error {
	fmt.Println("SignTraceID")
	encryptor := xhs.NewMiscEncrypt()
	traceID := encryptor.X_B3_TraceID()
	xrayTraceID := encryptor.X_Xray_TraceID(traceID)
	req.Header.Set("x-b3-traceid", traceID)
	req.Header.Set("x-xray-traceid", xrayTraceID)
	return nil
}
