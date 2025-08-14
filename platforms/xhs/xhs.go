package xhs

import (
	"context"
	"crawler-sdk/internal/crypto/xhs"
	"crawler-sdk/pkg/http"
	"encoding/base64"
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
	return &Client{
		client: http.NewClient(cookie),
	}
}

// GetUser 实现了获取用户信息的逻辑
func (c *Client) GetUser(ctx context.Context, userID string) (UserDetail, error) {
	// 注意：GetUserProfile 现在需要从 cookie 中获取 a1，但 NewClient 已经将 cookie 设置在请求头中。
	// 为了加密，我们仍然需要原始的 cookie 字符串。
	// 理想情况下，可以将 cookie 存储在 Client 结构体中。
	// 这里为了简化，我们直接在 GetUserProfile 中传递它。
	// 一个更好的重构是在 New 函数中解析 a1 并存储。
	// 但为了保持与 Python 版本的接口一致性，我们暂时在 GetUser 中传递。
	return c.getUserDetail(ctx, userID)
}

func SignXYS(c *resty.Client, req *resty.Request) error {
	xsc := xhs.XYS(req.URL, req.Header.Get("Cookie"))
	req.Header.Set("x-s", xsc)
	return nil
}

func SignXS(c *resty.Client, req *resty.Request) error {
	encryptor := xhs.NewXsEncrypt()
	cookie := req.Header.Get("Cookie")
	a1 := getCookieValue(cookie, "a1")
	ts := fmt.Sprintf("%d", time.Now().Unix()*1000)
	signedURL, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("URL解析失败: %w", err)
	}
	q := signedURL.Query()
	for k, v := range req.QueryParams {
		q.Set(k, v[0])
	}
	signedURL.RawQuery = q.Encode()

	xs, err := encryptor.EncryptXs(signedURL.String(), a1, ts, "xhs-pc-web")
	if err != nil {
		return fmt.Errorf("加密 'x-s' 失败: %w", err)
	}
	req.Header.Set("x-s", xs)
	req.Header.Set("x-t", ts)
	return nil
}

func SignXSC(c *resty.Client, req *resty.Request) error {
	encryptor := xhs.XscEncrypt{}
	cookie := req.Header.Get("Cookie")
	a1 := getCookieValue(cookie, "a1")
	x1 := getCookieValue(cookie, "x1")
	x4 := getCookieValue(cookie, "x4")
	b1 := getCookieValue(cookie, "b1")
	xsc, err := encryptor.EncryptXsc(req.Header.Get("x-s"), req.Header.Get("x-t"), "xhs-pc-web", a1, x1, x4, b1)
	if err != nil {
		return fmt.Errorf("加密 'x-sc' 失败: %w", err)
	}
	xscStr := base64.StdEncoding.EncodeToString(xsc)
	req.Header.Set("x-sc", xscStr)
	return nil
}

func SignTraceID(c *resty.Client, req *resty.Request) error {
	encryptor := xhs.NewMiscEncrypt()
	traceID := encryptor.X_B3_TraceID()
	xrayTraceID := encryptor.X_Xray_TraceID(traceID)
	req.Header.Set("x-b3-traceid", traceID)
	req.Header.Set("x-xray-traceid", xrayTraceID)
	return nil
}
