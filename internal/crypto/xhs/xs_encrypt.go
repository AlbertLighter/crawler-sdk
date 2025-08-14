package xhs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// XsEncrypt 提供了小红书 xs 相关的加密功能
type XsEncrypt struct {
	words    []uint32
	keyBytes []byte
	iv       []byte
}

// NewXsEncrypt 创建 XsEncrypt 实例
func NewXsEncrypt() *XsEncrypt {
	words := []uint32{929260340, 1633971297, 895580464, 925905270}
	keyBytes := make([]byte, 16)
	for i, word := range words {
		keyBytes[i*4] = byte(word >> 24)
		keyBytes[i*4+1] = byte(word >> 16)
		keyBytes[i*4+2] = byte(word >> 8)
		keyBytes[i*4+3] = byte(word)
	}
	iv := []byte("4uzjr7mbsibcaldp")

	return &XsEncrypt{
		words:    words,
		keyBytes: keyBytes,
		iv:       iv,
	}
}

// pkcs7Padding 实现了 PKCS7 填充
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// EncryptMD5 根据传入的url生成MD5摘要
func (x *XsEncrypt) EncryptMD5(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

// EncryptText 根据传入的text生成AES加密后的内容，并将其转为base64编码
func (x *XsEncrypt) EncryptText(text string) (string, error) {
	textEncoded := base64.StdEncoding.EncodeToString([]byte(text))

	block, err := aes.NewCipher(x.keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建AES cipher失败: %w", err)
	}

	paddedText := pkcs7Padding([]byte(textEncoded), aes.BlockSize)
	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, x.iv)
	mode.CryptBlocks(ciphertext, paddedText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Base64ToHex 把加密后的payload转为16进制
func (x *XsEncrypt) Base64ToHex(encodedData string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}
	return hex.EncodeToString(decodedData), nil
}

// EncryptPayload 把小红书加密参数payload转16进制 再使用base64编码
func (x *XsEncrypt) EncryptPayload(payload, platform string) (string, error) {
	hexPayload, err := x.Base64ToHex(payload)
	if err != nil {
		return "", fmt.Errorf("payload转16进制失败: %w", err)
	}

	obj := map[string]string{
		"signSvn":     "56",
		"signType":    "x2",
		"appID":       platform,
		"signVersion": "1",
		"payload":     hexPayload,
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// EncryptXs 将传入的参数加密为小红书的xs
func (x *XsEncrypt) EncryptXs(url, a1, ts, platform string) (string, error) {
	md5URL := x.EncryptMD5("url=" + url)
	text := fmt.Sprintf("x1=%s;x2=0|0|0|1|0|0|1|0|0|0|1|0|0|0|0|1|0|0|0;x3=%s;x4=%s;", md5URL, a1, ts)

	encryptedText, err := x.EncryptText(text)
	if err != nil {
		return "", fmt.Errorf("加密文本失败: %w", err)
	}

	encryptedPayload, err := x.EncryptPayload(encryptedText, platform)
	if err != nil {
		return "", fmt.Errorf("加密payload失败: %w", err)
	}

	return "XYW_" + encryptedPayload, nil
}

// EncryptSign 小红书验证码签名
func (x *XsEncrypt) EncryptSign(ts string, payload map[string]interface{}) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("JSON编码payload失败: %w", err)
	}
	url := fmt.Sprintf("%stest/api/redcaptcha/v2/captcha/register%s", ts, string(payloadBytes))

	md5URL := x.EncryptMD5(url)
	md5ASCII := []byte(md5URL)

	var result strings.Builder
	for i := 0; i < len(md5ASCII); i += 3 {
		u := int(md5ASCII[i])
		c := 0
		s := 0

		if i+1 < len(md5ASCII) {
			c = int(md5ASCII[i+1])
		}
		if i+2 < len(md5ASCII) {
			s = int(md5ASCII[i+2])
		}

		l := u >> 2
		f := ((u & 3) << 4) | (c >> 4)
		p := ((c & 15) << 2) | (s >> 6)
		d := s & 63

		result.WriteByte(XN[l])
		result.WriteByte(XN[f])

		if p < 64 {
			result.WriteByte(XN[p])
		} else {
			result.WriteByte(XN64)
		}
		if d < 64 {
			result.WriteByte(XN[d])
		} else {
			result.WriteByte(XN64)
		}
	}
	return result.String(), nil
}
