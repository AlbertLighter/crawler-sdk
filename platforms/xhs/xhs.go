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

type XhsClient struct {
	client       *resty.Client
	UploadClient *resty.Client
}

// New 创建一个新的小红书客户端实例
// cookie: 用于身份验证的小红书 cookie 字符串
func New(cookie string) *XhsClient {
	c := &XhsClient{
		client:       http.NewClient(cookie),
		UploadClient: resty.New(),
	}
	c.client.AddRequestMiddleware(headers)
	c.client.AddRequestMiddleware(SignXYS)
	// c.client.AddRequestMiddleware(SignXS)
	c.client.AddRequestMiddleware(SignXSC)
	c.client.AddRequestMiddleware(SignTraceID)
	c.UploadClient.AddRequestMiddleware(QSign)
	c.client.SetProxy("http://127.0.0.1:8888")
	c.UploadClient.SetProxy("http://127.0.0.1:8888")
	return c
}

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
	cookie := c.Header().Get("Cookie")
	u, err := url.Parse(req.URL)
	if err != nil {
		return err
	}
	fmt.Print(u)
	a1 := getCookieValue(cookie, "a1")
	// xsc := xhs.XYS(req.URL, a1)
	fmt.Println(req.Method)
	fmt.Println(u.Path)
	fmt.Println(a1)
	fmt.Println("xhs-pc-web")
	fmt.Println(req.QueryParams)
	xsc := xhs.Sign(req.Method, u.Path, a1, "xhs-pc-web", req.QueryParams)
	req.Header.Set("X-s", xsc)
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())
	req.Header.Set("X-t", ts)

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
	x1 := "4.2.4"
	x4 := "0.79.4"
	b1 := `I38rHdgsjopgIvesdVwgIC+oIELmBZ5e3VwXLgFTIxS3bqwErFeexd0ekncAzMFYnqthIhJeDnMDKutRI3KsYorWHPtGrbV0P9WfIi/eWc6eYqtyQApPI37ekmR6QL+5Ii6sdneeSfqYHqwl2qt5B0DBIx+PGDi/sVtkIxdsxuwr4qtiIhuaIE3e3LV0I3VTIC7e0utl2ADmsLveDSKsSPw5IEvsiVtJOqw8BuwfPpdeTFWOIx4TIiu6ZPwrPut5IvlaLbgs3qtxIxes1VwHIkumIkIyejgsY/WTge7eSqte/D7sDcpipedeYrDtIC6eDVw2IENsSqtlnlSuNjVtIvoekqt3cZ7sVo4gIESyIhEG+9DUIvzy4I8OIic7ZPwAIviu4o/sDLds6PwVIC7eSd7ej/Y4IEve1SiMtVwUIids3s/sxZNeiVtbcUeeYVwvIvkazA0eSVwhLfKsfPwoIxltIxZSouwOgVwpsr4heU/e6LveDPwFIvgs1ros1DZiIi7sjbos3grFIE0e3PtvIibROqwOOqthIxesxPw7IhAeVPthIh/sYqtSGqwymPwDIvIkI3It4aGS4Y/eiutjIimrIEOsSVtzBoFM/9vej9ZvIiENGutzrutlIvve3PtUOpKe1Y6s3LMoIh7sVd0siPtPLuwwIveeSPwRIiNeksrLI37eD9KeWVtRI37sxuwU4eQNIEH+Iv7sxM6ex7vsYDosSPtzIkL1IE4RaPtLICrYIEgei/iEGUKsWVtbIEPZzVwuwS7eVI/sfPta/Pt2IxPAIERLeVtLmVwBIxvsYr0e3I==`
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

func QSign(c *resty.Client, req *resty.Request) error {
	cookie := c.Header().Get("Cookie")
	a1 := getCookieValue(cookie, "a1")
	q, err := xhs.GetQSignAuth(xhs.QSignAuthOptions{
		SecretId:  a1,
		SecretKey: a1,
		Method:    req.Method,
		Pathname:  req.URL,
		Query:     map[string]string{},
		Headers:   map[string]string{},
	})
	if err != nil {
		return fmt.Errorf("获取QSignAuth失败: %w", err)
	}
	req.Header.Set("Authorization", q)
	return nil
}
