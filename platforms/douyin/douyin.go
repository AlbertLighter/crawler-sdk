package douyin

import (
	"crawler-sdk/internal/crypto/dy/ab"
	"crawler-sdk/internal/crypto/dy/auth"
	"crawler-sdk/pkg/http"
	"math/rand"

	"fmt"
	"time"

	"resty.dev/v3"
)

// AuthDetails holds the temporary credentials for uploading.
// It's derived from the /web/api/media/upload/auth/v5/ endpoint response.

// DyClient 是抖音平台的客户端
type DyClient struct {
	client       *resty.Client
	imagexClient *resty.Client
	uploadClient *resty.Client
}

// New 创建一个新的抖音客户端实例
// cookie: 用于身份验证的抖音 cookie 字符串
func New(cookie string) *DyClient {
	c := &DyClient{
		client:       http.NewClient(cookie),
		imagexClient: http.NewClient(""),
		uploadClient: http.NewClient(""),
	}
	c.client.AddRequestMiddleware(headers)
	c.client.AddRequestMiddleware(msToken)
	c.client.AddRequestMiddleware(absign)
	c.imagexClient.AddRequestMiddleware(uploadSign)
	c.imagexClient.AddRequestMiddleware(uploadHeaders)
	// c.uploadClient.AddRequestMiddleware(uploadSign)
	// c.uploadClient.AddRequestMiddleware(uploadHeaders)
	return c
}

func msToken(c *resty.Client, req *resty.Request) error {
	msToken := getMsToken(120)
	req.SetQueryParam("msToken", msToken)
	return nil
}

func absign(c *resty.Client, req *resty.Request) error {
	fmt.Println("a_bogus")
	queryStr := req.QueryParams.Encode()
	sign, err := ab.SignDetail(queryStr, c.Header()["User-Agent"][0])
	if err != nil {
		return err
	}
	req.SetQueryParam("a_bogus", sign)
	return nil
}

// def get_ms_token(randomlength=120):
// """
// 根据传入长度产生随机字符串
// """
// random_str = ''
// base_str = 'ABCDEFGHIGKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz0123456789='
// length = len(base_str) - 1
// for _ in range(randomlength):
// 	random_str += base_str[random.randint(0, length)]
// return random_str

func getMsToken(randomlength int) string {
	random_str := ""
	base_str := "ABCDEFGHIGKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz0123456789="
	length := len(base_str) - 1
	for range randomlength {
		random_str += string(base_str[rand.Intn(length)])
	}
	return random_str
}

func headers(c *resty.Client, req *resty.Request) error {
	fmt.Println("headers")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-ch-ua", `"Not;A=Brand";v="99", "Microsoft Edge";v="139", "Chromium";v="139"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	return nil
}

func uploadHeaders(c *resty.Client, req *resty.Request) error {
	req.Header.Set("Host", "imagex.bytedanceapi.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")
	req.Header.Set("sec-ch-ua", `"Not;A=Brand";v="99", "Microsoft Edge";v="139", "Chromium";v="139"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://creator.douyin.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://creator.douyin.com/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	return nil
}

func uploadSign(c *resty.Client, req *resty.Request) error {
	fmt.Println("uploadSign")
	a := req.Context().Value("auth").(*AuthDetails)
	signer := auth.NewSigner(req, "imagex", "cn-north-1", false)
	credentials := auth.Credentials{
		AccessKeyID:     a.AccessKeyID,
		SecretAccessKey: a.SecretAccessKey,
		SessionToken:    a.SessionToken,
	}
	signer.AddAuthorization(credentials, time.Now())
	return nil
}
