package bd

import (
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"golang.org/x/crypto/hkdf"
)

// security-sdk/s_sdk_crypt_sdk	{"data": "{“ec_privateKey“:“-----BEGIN PRIVATE KEY-----\\r\\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg79y2ijwRFYYWvOTG\\r\\nQoYsueKWR0HKb28VYt6aBKgJEsChRANCAAT1OZAO2ROuib5C4tyNNpMhTaBnEuuK\\r\\ni0HQz9ms0EARQeENlNRVyK1UpWyWjniFr0pCd7UvYz/EJEkVGSbMj6Uq\\r\\n-----END PRIVATE KEY-----\\r\\n“,“ec_publicKey“:“-----BEGIN PUBLIC KEY-----\\r\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE9TmQDtkTrom+QuLcjTaTIU2gZxLr\\r\\niotB0M/ZrNBAEUHhDZTUVcitVKVslo54ha9KQne1L2M/xCRJFRkmzI+lKg==\\r\\n-----END PUBLIC KEY-----\\r\\n“,“ec_csr“:“-----BEGIN CERTIFICATE REQUEST-----\\r\\nMIIBDjCBtQIBADAnMQswCQYDVQQGEwJDTjEYMBYGA1UEAwwPYmRfdGlja2V0X2d1\\r\\nYXJkMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE9TmQDtkTrom+QuLcjTaTIU2g\\r\\nZxLriotB0M/ZrNBAEUHhDZTUVcitVKVslo54ha9KQne1L2M/xCRJFRkmzI+lKqAs\\r\\nMCoGCSqGSIb3DQEJDjEdMBswGQYDVR0RBBIwEIIOd3d3LmRvdXlpbi5jb20wCgYI\\r\\nKoZIzj0EAwIDSAAwRQIhANHodtkMmyed5mLRjBydvohgyeT0SdOuosDRsFTe+gJy\\r\\nAiAU8rajUFBDRTwybutyiwR+5HPdtYU+Pcq6lDZX/8fxag==\\r\\n-----END CERTIFICATE REQUEST-----\\r\\n“}"}

// security-sdk/s_sdk_sign_data_key/web_protect	{"data":"{\"ticket\":\"hash.fvsLW/x06mCbBZxVS/RFz/Ly8/P9TdEBGIVKmPHY3mw=\",\"ts_sign\":\"ts.2.9325933b8bf09d548974044a575a30e44f819f094a033643db648fa7ad7fd41bc4fbe87d2319cf05318624ceda14911ca406dedbebeddb2e30fce8d4fa02575d\",\"client_cert\":\"pub.BPU5kA7ZE66JvkLi3I02kyFNoGcS64qLQdDP2azQQBFB4Q2U1FXIrVSlbJaOeIWvSkJ3tS9jP8QkSRUZJsyPpSo=\",\"log_id\":\"202507301814285E344115E4D37136DFF1\",\"create_time\":1753870469}"}

// security-sdk/s_sdk_server_cert_key	{"cert":"-----BEGIN CERTIFICATE-----\nMIIEfTCCBCKgAwIBAgIUXWdS2tzmSoewCWfKFyiWMrJqs/0wCgYIKoZIzj0EAwIw\nMTELMAkGA1UEBhMCQ04xIjAgBgNVBAMMGXRpY2tldF9ndWFyZF9jYV9lY2RzYV8y\nNTYwIBcNMjIxMTE4MDUyMDA2WhgPMjA2OTEyMzExNjAwMDBaMCQxCzAJBgNVBAYT\nAkNOMRUwEwYDVQQDEwxlY2llcy1zZXJ2ZXIwWTATBgcqhkjOPQIBBggqhkjOPQMB\nBwNCAASE2llDPlfc8Rq+5J5HXhg4edFjPnCF3Ua7JBoiE/foP9m7L5ELIcvxCgEx\naRCHbQ8kCCK/ArZ4FX/qCobZAkToo4IDITCCAx0wDgYDVR0PAQH/BAQDAgWgMDEG\nA1UdJQQqMCgGCCsGAQUFBwMBBggrBgEFBQcDAgYIKwYBBQUHAwMGCCsGAQUFBwME\nMCkGA1UdDgQiBCABydxqGrVEHhtkCWTb/vicGpDZPFPDxv82wiuywUlkBDArBgNV\nHSMEJDAigCAypWfqjmRIEo3MTk1Ae3MUm0dtU3qk0YDXeZSXeyJHgzCCAZQGCCsG\nAQUFBwEBBIIBhjCCAYIwRgYIKwYBBQUHMAGGOmh0dHA6Ly9uZXh1cy1wcm9kdWN0\naW9uLmJ5dGVkYW5jZS5jb20vYXBpL2NlcnRpZmljYXRlL29jc3AwRgYIKwYBBQUH\nMAGGOmh0dHA6Ly9uZXh1cy1wcm9kdWN0aW9uLmJ5dGVkYW5jZS5uZXQvYXBpL2Nl\ncnRpZmljYXRlL29jc3AwdwYIKwYBBQUHMAKGa2h0dHA6Ly9uZXh1cy1wcm9kdWN0\naW9uLmJ5dGVkYW5jZS5jb20vYXBpL2NlcnRpZmljYXRlL2Rvd25sb2FkLzQ4RjlD\nMEU3QjBDNUE3MDVCOTgyQkU1NTE3MDVGNjQ1QzhDODc4QTguY3J0MHcGCCsGAQUF\nBzAChmtodHRwOi8vbmV4dXMtcHJvZHVjdGlvbi5ieXRlZGFuY2UubmV0L2FwaS9j\nZXJ0aWZpY2F0ZS9kb3dubG9hZC80OEY5QzBFN0IwQzVBNzA1Qjk4MkJFNTUxNzA1\nRjY0NUM4Qzg3OEE4LmNydDCB5wYDVR0fBIHfMIHcMGygaqBohmZodHRwOi8vbmV4\ndXMtcHJvZHVjdGlvbi5ieXRlZGFuY2UuY29tL2FwaS9jZXJ0aWZpY2F0ZS9jcmwv\nNDhGOUMwRTdCMEM1QTcwNUI5ODJCRTU1MTcwNUY2NDVDOEM4NzhBOC5jcmwwbKBq\noGiGZmh0dHA6Ly9uZXh1cy1wcm9kdWN0aW9uLmJ5dGVkYW5jZS5uZXQvYXBpL2Nl\ncnRpZmljYXRlL2NybC80OEY5QzBFN0IwQzVBNzA1Qjk4MkJFNTUxNzA1RjY0NUM4\nQzg3OEE4LmNybDAKBggqhkjOPQQDAgNJADBGAiEAqMjT5ADMdGMeaImoJK4J9jzE\nLqZ573rNjsT3k14pK50CIQCLpWHVKWi71qqqrMjiSDvUhpyO1DpTPRHlavPRuaNm\nww==\n-----END CERTIFICATE-----","sn":"533240336124694022040808462028007165443034493949","createdTime":1755846638240}

type BDSigner struct {
	// ts_sign
	TsSign string `json:"ts_sign"` //security-sdk/s_sdk_crypt_sdk
	// ec_privateKey
	EcPrivateKey string `json:"ec_private_key"` //security-sdk/s_sdk_sign_data_key/web_protect
	// ec_publicKey
	EcPublicKey string `json:"ec_public_key"` //security-sdk/s_sdk_sign_data_key/web_protect
	// cert
	Cert string `json:"cert"` //security-sdk/s_sdk_server_cert_key
	// ticket
	Ticket string `json:"ticket"`

	derivedKey   []byte
	reePublicKey string
}

// {"ts_sign":"ts.2.9325933b8bf09d548974044a575a30e44f819f094a033643db648fa7ad7fd41bc4fbe87d2319cf05318624ceda14911ca406dedbebeddb2e30fce8d4fa02575d","req_content":"ticket,path,timestamp","req_sign":"vr7Kg8mtcaTqTmJdRa9nZm5RE/0eSGq+IORFtq3WljU=","timestamp":1756023357}

// bd-ticket-guard-client-data
type BDTicketGuardClientData struct {
	TsSign     string `json:"ts_sign"`
	ReqContent string `json:"req_content"`
	ReqSign    string `json:"req_sign"`
	Timestamp  int64  `json:"timestamp"`
}

func NewBDSigner(tsSign string, ticket string, ecPrivateKey string, ecPublicKey string, cert string) (*BDSigner, error) {
	bd := &BDSigner{
		TsSign:       tsSign,
		Ticket:       ticket,
		EcPrivateKey: ecPrivateKey,
		EcPublicKey:  ecPublicKey,
		Cert:         cert,
	}
	var err error
	bd.derivedKey, err = bd.DeriveECDHKey(ecPrivateKey, cert)
	if err != nil {
		return nil, err
	}
	bd.reePublicKey, err = GetPubKeyBase64(bd.EcPrivateKey, bd.EcPublicKey)
	if err != nil {
		return nil, err
	}
	return bd, nil
}

// "ticket=hash.fvsLW/x06mCbBZxVS/RFz/Ly8/P9TdEBGIVKmPHY3mw=&path=/aweme/janus/creator/comment/aweme/v1/comment/publish/&timestamp=1756024739"
func (b *BDSigner) BDSign(path string) (string, error) {
	t := time.Now().Unix()
	reqContent := fmt.Sprintf("ticket=%s&path=%s&timestamp=%d", b.Ticket, path, t)
	reqSign := b.SignWithHmacSha256([]byte(reqContent))
	d := &BDTicketGuardClientData{
		TsSign:     b.TsSign,
		ReqContent: "ticket,path,timestamp",
		ReqSign:    base64.StdEncoding.EncodeToString([]byte(reqSign)),
		Timestamp:  t,
	}
	jsonData, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// DeriveECDHKey 使用 ECDH 和 HKDF 派生共享密钥。
// privateKeyPath: 本地私钥 PEM 文件的路径
// peerCertPath: 对端证书 PEM 文件的路径
// 返回派生出的 32 字节密钥
func (b *BDSigner) DeriveECDHKey(privateKeyStr string, cert string) ([]byte, error) {
	// 1. 加载本地私钥
	privKeyBlock, _ := pem.Decode([]byte(privateKeyStr))
	if privKeyBlock == nil {
		return nil, fmt.Errorf("无法解码 PEM 格式的私钥")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	ecdhPrivKey, ok := privateKey.(*ecdh.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥不是有效的 ECDH 私钥")
	}

	// 2. 加载对端的证书并提取公钥
	certBlock, _ := pem.Decode([]byte(cert))
	if certBlock == nil {
		return nil, fmt.Errorf("无法解码 PEM 格式的证书")
	}
	peerCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %w", err)
	}

	ecdhPubKey, ok := peerCert.PublicKey.(*ecdh.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥不是有效的 ECDH 公钥")
	}

	// 3. 执行 ECDH 获取共享密钥
	sharedSecret, err := ecdhPrivKey.ECDH(ecdhPubKey)
	if err != nil {
		return nil, fmt.Errorf("ECDH 密钥交换失败: %w", err)
	}

	// 4. 使用 HKDF 派生最终密钥
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, nil)
	derivedKey := make([]byte, 32)
	_, err = hkdfReader.Read(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("HKDF 密钥派生失败: %w", err)
	}

	return derivedKey, nil
}

// SignWithHmacSha256 使用 HMAC-SHA256 算法和给定的密钥对消息进行签名。
// key: 用于签名的密钥
// message: 需要签名的消息
// 返回签名后的十六进制字符串
func (b *BDSigner) SignWithHmacSha256(message []byte) string {
	// 创建一个新的 HMAC hash，使用 SHA256
	mac := hmac.New(sha256.New, b.derivedKey)
	// 写入待签名的消息
	mac.Write(message)
	// 计算并返回十六进制编码的签名
	return hex.EncodeToString(mac.Sum(nil))
}
