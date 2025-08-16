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
	"net/url"
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
	params := url.Values{}
	params.Set("aid", "1128")
	params.Set("support_h265", "1")
	params.Set("msToken", "E-qrN5_3s_GdJMz0_nDQm3nwc-Y4lyv0AAZPLCxZkofS0bz-I7dnsoY4_kVfOvQ5mvIHR6yYv6fJqx9OmOkaqd1Hju98GuCPs00pZecchjGT2n5fCx3cJxTD2sozSKRsZNCOLKfn6--0rMR9yexL9ml_hVToYy111jzJzOFeBNwc")
	params.Set("a_bogus", "Qv4VhFtEQx85OpFGuKEtC31UEWxlNP8yqlTQbzLn9PxKOZUGDZHakcGeGxztf7uxnYBVkKVHsfsAbxxbTUkzZA9pzmkfSNt6jzVInX8o01qDbzvsErjDSL6FoXsc8bGulQ5yiAXfMUt72xO-NrdD%2Fp-Hy%2FbF5QmkQrQRk%2FzGOoG11zyAE1c-PptkihiKUenJ")

	resp, err := c.client.R().
		SetHeaders(map[string]string{
			"Referer": "https://creator.douyin.com/creator-micro/content/post",
		}).
		SetResult(&UploadAuthResponse{}).
		Get(fmt.Sprintf("%s/web/api/media/upload/auth/v5/?%s", douyinCreatorURL, params.Encode()))

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
	reqURL := fmt.Sprintf("%s/?Action=ApplyImageUpload&Version=2018-08-01&ServiceId=jm8ajry58r", imagexURL)
	req, _ := http.NewRequest("GET", reqURL, nil)

	// AWS v4 签名
	headers := c.signRequest(ctx, req, auth, "")

	resp, err := c.client.R().
		SetHeaders(headers).
		SetResult(&ApplyUploadResponse{}).
		Get(reqURL)

	if err != nil {
		return nil, err
	}
	return resp.Result().(*ApplyUploadResponse), nil
}

// uploadFileToTOS 上传文件到TOS
func (c *DyClient) uploadFileToTOS(ctx context.Context, applyResp *ApplyUploadResponse, filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	crc32q := crc32.MakeTable(0xedb88320)
	crc := crc32.Checksum(fileBytes, crc32q)
	crc32Str := fmt.Sprintf("%x", crc)

	uploadURL := "https://" + applyResp.Result.UploadAddress.UploadHosts[0] + "/" + applyResp.Result.UploadAddress.StoreInfos[0].StoreUri

	resp, err := c.client.R().
		SetHeaders(map[string]string{
			"Content-CRC32": crc32Str,
			"Content-Type":  "application/octet-stream",
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

	resp, err := c.client.R().
		SetHeaders(headers).
		SetBody(payloadBytes).
		SetResult(&CommitUploadResponse{}).
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

// PublishImage 发布图文
func (c *DyClient) PublishImage(ctx context.Context, filePath, title, description string) (*CommitUploadResponse, error) {
	// 1. 获取上传凭证
	auth, err := c.getUploadAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("step 1: get upload auth failed: %w", err)
	}

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

	// 5. 发布作品 (此步骤需要分析发布接口)
	// imageUri := commitResp.Result.PluginResult[0].ImageUri
	// err = c.createPost(imageUri, title, description)
	// if err != nil {
	// 	 return nil, fmt.Errorf("step 5: create post failed: %w", err)
	// }
	fmt.Println("Image uploaded successfully! URI:", commitResp.Result.PluginResult[0].ImageUri)
	fmt.Println("Next step is to call the actual publish API which is not included in the log.")

	return commitResp, nil
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
