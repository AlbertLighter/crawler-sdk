package douyin

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"os"
	"time"
)

const (
	douyinCreatorURL = "https://creator.douyin.com"
	imagexURL        = "https://imagex.bytedanceapi.com"
)

// UploadAuthResponse 获取上传凭证的响应
type UploadAuthResponse struct {
	Ak          string `json:"ak"`
	Auth        string `json:"auth"`
	ClientIP    string `json:"client_ip"`
	ExpiredTime string `json:"ExpiredTime"`
	CurrentTime string `json:"CurrentTime"`
	StatusCode  int    `json:"status_code"`
}

// AuthDetails 解析后的凭证详情
type AuthDetails struct {
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	ExpiredTime     string `json:"ExpiredTime"`
	CurrentTime     string `json:"CurrentTime"`
}

// ApplyUploadResponse 申请上传的响应
type ApplyUploadResponse struct {
	ResponseMetadata struct {
		RequestID string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"cn-north-1"`
	} `json:"ResponseMetadata"`
	Result struct {
		UploadAddress struct {
			StoreInfos  []StoreInfo `json:"StoreInfos"`
			UploadHosts []string    `json:"UploadHosts"`
			SessionKey  string      `json:"SessionKey"`
		} `json:"UploadAddress"`
	} `json:"Result"`
}

// StoreInfo 存储信息
type StoreInfo struct {
	StoreUri string `json:"StoreUri"`
	Auth     string `json:"Auth"`
}

// CommitUploadResponse 确认上传的响应
type CommitUploadResponse struct {
	ResponseMetadata struct {
		RequestID string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"cn-north-1"`
	} `json:"ResponseMetadata"`
	Result struct {
		Results      []CommitResult `json:"Results"`
		RequestID    string         `json:"RequestId"`
		PluginResult []PluginResult `json:"PluginResult"`
	} `json:"Result"`
}

// CommitResult 确认上传的结果
type CommitResult struct {
	Uri       string `json:"Uri"`
	UriStatus int    `json:"UriStatus"`
}

// PluginResult 插件结果
type PluginResult struct {
	FileName    string `json:"FileName"`
	ImageUri    string `json:"ImageUri"`
	ImageWidth  int    `json:"ImageWidth"`
	ImageHeight int    `json:"ImageHeight"`
	ImageMd5    string `json:"ImageMd5"`
	ImageFormat string `json:"ImageFormat"`
	ImageSize   int    `json:"ImageSize"`
}

// getUploadAuth 获取上传凭证
func (c *DyClient) getUploadAuth(ctx context.Context) (*AuthDetails, error) {
	// TODO: msToken 和 a_bogus 参数需要动态生成
	// cookie_enabled	true
	// screen_width	2048
	// screen_height	1152
	// browser_language	zh-CN
	// browser_platform	Win32
	// browser_name	Mozilla
	// browser_version	5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0
	// browser_online	true
	// timezone_name	Asia/Shanghai
	// aid	1128
	// support_h265	1
	// msToken	E-qrN5_3s_GdJMz0_nDQm3nwc-Y4lyv0AAZPLCxZkofS0bz-I7dnsoY4_kVfOvQ5mvIHR6yYv6fJqx9OmOkaqd1Hju98GuCPs00pZecchjGT2n5fCx3cJxTD2sozSKRsZNCOLKfn6--0rMR9yexL9ml_hVToYy111jzJzOFeBNwc
	// a_bogus	Qv4VhFtEQx85OpFGuKEtC31UEWxlNP8yqlTQbzLn9PxKOZUGDZHakcGeGxztf7uxnYBVkKVHsfsAbxxbTUkzZA9pzmkfSNt6jzVInX8o01qDbzvsErjDSL6FoXsc8bGulQ5yiAXfMUt72xO-NrdD/p-Hy/bF5QmkQrQRk/zGOoG11zyAE1c-PptkihiKUenJ
	resp, err := c.client.R().
		SetHeaders(map[string]string{
			"Referer": "https://creator.douyin.com/creator-micro/content/post",
		}).
		SetResult(&UploadAuthResponse{}).
		SetQueryParams(map[string]string{
			"aid":             "1128",
			"support_h265":    "1",
			"cookie_enabled":  "true",
			"screen_width":    "2048",
			"screen_height":   "1152",
			"browser_language": "zh-CN",
			"browser_platform": "Win32",
			"browser_name":     "Mozilla",
			"browser_version":  "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0",
			"browser_online":   "true",
			"timezone_name":    "Asia/Shanghai",
		}).
		Get(fmt.Sprintf("%s/web/api/media/upload/auth/v5/", douyinCreatorURL))

	if err != nil {
		return nil, err
	}

	authResp := resp.Result().(*UploadAuthResponse)
	if authResp.StatusCode != 0 {
		return nil, fmt.Errorf("get upload auth failed, status code: %d", authResp.StatusCode)
	}

	var authDetails AuthDetails
	err = json.Unmarshal([]byte(authResp.Auth), &authDetails)
	if err != nil {
		return nil, err
	}

	return &authDetails, nil
}

// applyImageUpload 申请图片上传
func (c *DyClient) applyImageUpload(ctx context.Context, auth *AuthDetails) (*ApplyUploadResponse, error) {
	reqURL := fmt.Sprintf("%s/?Action=ApplyImageUpload&Version=2018-08-01&ServiceId=jm8ajry58r&app_id=2906&user_id=&s=p9t685goxl", imagexURL)
	// req, _ := http.NewRequest("GET", reqURL, nil)

	// AWS v4 签名
	// headers := c.signRequest(ctx, req, auth, "")

	resp, err := c.imagexClient.R().
		SetResult(&ApplyUploadResponse{}).
		SetContext(ctx).
		Get(reqURL)

	if err != nil {
		return nil, err
	}
	return resp.Result().(*ApplyUploadResponse), nil
}

// POST https://tos-d-x-hl.snssdk.com/upload/v1/tos-cn-i-jm8ajry58r/ebeb8e3b40ba4b448bddee77b52d0179 HTTP/1.1
// Host: tos-d-x-hl.snssdk.com
// Connection: keep-alive
// Content-Length: 16368
// sec-ch-ua-platform: "Windows"
// Authorization: SpaceKey/jm8ajry58r/1/:version:v2:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUzNDQwNDgsInNpZ25hdHVyZUluZm8iOnsiYWNjZXNzS2V5IjoiZmFrZV9hY2Nlc3Nfa2V5IiwiYnVja2V0IjoidG9zLWNuLWktam04YWpyeTU4ciIsImV4cGlyZSI6MTc1NTM0NDA0OCwiZmlsZUluZm9zIjpbeyJvaWRLZXkiOiJlYmViOGUzYjQwYmE0YjQ0OGJkZGVlNzdiNTJkMDE3OSIsImZpbGVUeXBlIjoiMSJ9XSwiZXh0cmEiOnsiYWNjb3VudF9wcm9kdWN0IjoiaW1hZ2V4IiwiYmxvY2tfbW9kZSI6IiIsImNvbnRlbnRfdHlwZV9ibG9jayI6IntcIm1pbWVfcGN0XCI6MCxcIm1vZGVcIjowLFwibWltZV9saXN0XCI6bnVsbCxcImNvbmZsaWN0X2Jsb2NrXCI6ZmFsc2V9IiwiZW5jcnlwdF9hbGdvIjoiIiwiZW5jcnlwdF9rZXkiOiIiLCJzcGFjZSI6ImptOGFqcnk1OHIiLCJ0b3NfbWV0YSI6IntcIlVTRVJfSURcIjpcIjE4NzMxNTYwMjUwOTEyNTZcIn0ifX19.YxVzRhg_E4LlJAcE_Hez7-RuC7iLk5Tw_RrReAKbTdw
// sec-ch-ua: "Not;A=Brand";v="99", "Microsoft Edge";v="139", "Chromium";v="139"
// Content-CRC32: fc5f0689
// sec-ch-ua-mobile: ?0
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0
// X-Storage-U:
// Content-Type: application/octet-stream
// Content-Disposition: attachment; filename="undefined"
// Accept: */*
// Origin: https://creator.douyin.com
// Sec-Fetch-Site: cross-site
// Sec-Fetch-Mode: cors
// Sec-Fetch-Dest: empty
// Referer: https://creator.douyin.com/
// Accept-Encoding: gzip, deflate, br, zstd
// Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6

// uploadFileToTOS 上传文件到TOS
func (c *DyClient) uploadFileToTOS(ctx context.Context, applyResp *ApplyUploadResponse, filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	crc32q := crc32.MakeTable(0xedb88320)
	crc := crc32.Checksum(fileBytes, crc32q)
	crc32Str := fmt.Sprintf("%x", crc)

	uploadURL := "https://" + applyResp.Result.UploadAddress.UploadHosts[0] + "/upload/v1/" + applyResp.Result.UploadAddress.StoreInfos[0].StoreUri

	resp, err := c.uploadClient.R().
		SetHeaders(map[string]string{
			"Content-CRC32": crc32Str,
			"Content-Type":  "application/octet-stream",
			"Authorization": applyResp.Result.UploadAddress.StoreInfos[0].Auth,
		}).
		SetBody(fileBytes).
		Post(uploadURL)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("upload file failed, status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// commitImageUpload 确认图片上传
func (c *DyClient) commitImageUpload(ctx context.Context, auth *AuthDetails, applyResp *ApplyUploadResponse) (*CommitUploadResponse, error) {
	reqURL := fmt.Sprintf("%s/?Action=CommitImageUpload&Version=2018-08-01&ServiceId=jm8ajry58r", imagexURL)

	payload := map[string]string{
		"SessionKey": applyResp.Result.UploadAddress.SessionKey,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", reqURL, bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	// AWS v4 签名
	headers := c.signRequest(ctx, req, auth, string(payloadBytes))

	resp, err := c.imagexClient.R().
		SetHeaders(headers).
		SetBody(payloadBytes).
		SetResult(&CommitUploadResponse{}).
		SetContext(ctx).
		Post(reqURL)

	if err != nil {
		return nil, err
	}

	commitResp := resp.Result().(*CommitUploadResponse)
	if len(commitResp.Result.Results) > 0 && commitResp.Result.Results[0].UriStatus != 2000 {
		return nil, fmt.Errorf("commit upload failed with status: %d", commitResp.Result.Results[0].UriStatus)
	}

	return commitResp, nil
}

// CreatePostRequest is the request body for creating a post
type CreatePostRequest struct {
	Item Item `json:"item"`
}

// Item represents the item details in the post
type Item struct {
	Common  Common  `json:"common"`
	Cover   Cover   `json:"cover"`
	Anchor  Anchor  `json:"anchor"`
	Declare Declare `json:"declare"`
}

// Common contains common details of the post
type Common struct {
	Text           string  `json:"text"`
	TextExtra      string  `json:"text_extra"`
	Activity       string  `json:"activity"`
	Challenges     string  `json:"challenges"`
	HashtagSource  string  `json:"hashtag_source"`
	Mentions       string  `json:"mentions"`
	VisibilityType int     `json:"visibility_type"`
	Download       int     `json:"download"`
	Timing         int     `json:"timing"`
	MediaType      int     `json:"media_type"`
	Images         []Image `json:"images"`
	CreationID     string  `json:"creation_id"`
}

// Image represents an image in the post
type Image struct {
	URI    string `json:"uri"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Cover represents the cover of the post
type Cover struct {
	Poster string `json:"poster"`
}

// Anchor is an empty struct for now
type Anchor struct{}

// Declare contains declaration info
type Declare struct {
	UserDeclareInfo string `json:"user_declare_info"`
}

// CreatePostResponse is the response from creating a post
type CreatePostResponse struct {
	Extra struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
	ItemID     string `json:"item_id"`
	StatusCode int    `json:"status_code"`
}

// createPost creates a new post (aweme)
func (c *DyClient) createPost(ctx context.Context, commitResp *CommitUploadResponse, title string) (*CreatePostResponse, error) {
	// TODO: msToken and a_bogus need to be generated dynamically
	imageInfo := commitResp.Result.PluginResult[0]

	// Generate a creation_id, it seems to be a timestamp-based unique id.
	creationID := fmt.Sprintf("gemini%d", time.Now().UnixNano()/1e6)

	textExtra := fmt.Sprintf(`[{"start":0,"end":%d,"hashtag_id":0,"hashtag_name":"","type":7}]`, len(title))

	payload := CreatePostRequest{
		Item: Item{
			Common: Common{
				Text:           title,
				TextExtra:      textExtra,
				Activity:       "[]",
				Challenges:     "[]",
				HashtagSource:  "",
				Mentions:       "[]",
				VisibilityType: 0,
				Download:       1,
				Timing:         -1,
				MediaType:      2, // 2 for images
				Images: []Image{
					{
						URI:    imageInfo.ImageUri,
						Width:  imageInfo.ImageWidth,
						Height: imageInfo.ImageHeight,
					},
				},
				CreationID: creationID,
			},
			Cover: Cover{
				Poster: imageInfo.ImageUri,
			},
			Anchor:  Anchor{},
			Declare: Declare{
				UserDeclareInfo: "{}",
			},
		},
	}

	// Query params from the captured request
	queryParams := map[string]string{
		"read_aid":         "2906",
		"cookie_enabled":   "true",
		"screen_width":     "2048",
		"screen_height":    "1152",
		"browser_language": "zh-CN",
		"browser_platform": "Win32",
		"browser_name":     "Mozilla",
		"browser_version":  "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0",
		"browser_online":   "true",
		"timezone_name":    "Asia/Shanghai",
		"aid":              "1128",
		"support_h265":     "1",
		// "msToken": "...", // TODO: Needs to be generated dynamically
		// "a_bogus": "...", // TODO: Needs to be generated dynamically
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&CreatePostResponse{}).
		SetBody(payload).
		SetHeaders(map[string]string{
			"Referer":      "https://creator.douyin.com/creator-micro/content/post/image?default-tab=3&enter_from=publish_page&media_type=image&type=new",
			"Content-Type": "application/json",
		}).
		SetQueryParams(queryParams).
		Post(fmt.Sprintf("%s/web/api/media/aweme/create_v2/", douyinCreatorURL))

	if err != nil {
		return nil, err
	}

	createResp := resp.Result().(*CreatePostResponse)
	if createResp.StatusCode != 0 {
		return nil, fmt.Errorf("create post failed, status code: %d, item_id: %s", createResp.StatusCode, createResp.ItemID)
	}

	return createResp, nil
}

// PublishImage 发布图文
func (c *DyClient) PublishImage(ctx context.Context, filePath, title, description string) (*CreatePostResponse, error) {
	// 1. 获取上传凭证
	auth, err := c.getUploadAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("step 1: get upload auth failed: %w", err)
	}

	ctx = context.WithValue(ctx, "auth", auth)

	// 2. 申请上传
	applyResp, err := c.applyImageUpload(ctx, auth)
	if err != nil {
		return nil, fmt.Errorf("step 2: apply image upload failed: %w", err)
	}

	// 3. 上传文件
	err = c.uploadFileToTOS(ctx, applyResp, filePath)
	if err != nil {
		return nil, fmt.Errorf("step 3: upload file to TOS failed: %w", err)
	}

	// 4. 确认上传
	commitResp, err := c.commitImageUpload(ctx, auth, applyResp)
	if err != nil {
		return nil, fmt.Errorf("step 4: commit image upload failed: %w", err)
	}

	// 5. 发布作品
	createResp, err := c.createPost(ctx, commitResp, title)
	if err != nil {
		return nil, fmt.Errorf("step 5: create post failed: %w", err)
	}

	fmt.Println("Image published successfully! ItemID:", createResp.ItemID)

	return createResp, nil
}

// signRequest 为ImageX API请求进行AWS v4签名 (简化版)
func (c *DyClient) signRequest(ctx context.Context, req *http.Request, auth *AuthDetails, payload string) map[string]string {
	t := time.Now().UTC()
	amzDate := t.Format("20060102T150405Z")
	dateStamp := t.Format("20060102")

	// 1. 创建规范请求
	canonicalURI := req.URL.Path
	canonicalQueryString := req.URL.RawQuery
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-date:%s\n", req.URL.Host, amzDate)
	signedHeaders := "host;x-amz-date"

	payloadHash := sha256.Sum256([]byte(payload))
	payloadHashStr := hex.EncodeToString(payloadHash[:])

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHashStr)

	// 2. 创建待签字符串
	algorithm := "AWS4-HMAC-SHA256"
	credentialScope := fmt.Sprintf("%s/cn-north-1/imagex/aws4_request", dateStamp)
	canonicalRequestHash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s",
		algorithm, amzDate, credentialScope, hex.EncodeToString(canonicalRequestHash[:]))

	// 3. 计算签名
	signingKey := getSignatureKey(auth.SecretAccessKey, dateStamp, "cn-north-1", "imagex")
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// 4. 构建 Authorization 头部
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, auth.AccessKeyID, credentialScope, signedHeaders, signature)

	headers := map[string]string{
		"Authorization":        authorizationHeader,
		"X-Amz-Date":           amzDate,
		"X-Amz-Security-Token": auth.SessionToken,
		"Content-Type":         req.Header.Get("Content-Type"),
	}
	if payload != "" {
		headers["X-Amz-Content-Sha256"] = payloadHashStr
	}

	return headers
}

func hmacSHA256(key []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func getSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+key), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(regionName))
	kService := hmacSHA256(kRegion, []byte(serviceName))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}
