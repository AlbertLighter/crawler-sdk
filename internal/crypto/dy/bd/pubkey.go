package bd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// GetPubKeyBase64 从 PEM 格式的私钥或公钥中提取未压缩的公钥，并返回其 Base64 编码的字符串。
// privateKeyPEM: PEM 格式的私钥字符串
// publicKeyPEM: PEM 格式的公钥字符串
func GetPubKeyBase64(privateKeyPEM, publicKeyPEM string) (string, error) {
	var publicKey interface{}
	var err error

	// 优先使用公钥
	if publicKeyPEM != "" {
		pubKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
		if pubKeyBlock == nil {
			return "", fmt.Errorf("无法解码 PEM 格式的公钥")
		}
		publicKey, err = x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
		if err != nil {
			return "", fmt.Errorf("解析公钥失败: %w", err)
		}
	} else if privateKeyPEM != "" {
		privKeyBlock, _ := pem.Decode([]byte(privateKeyPEM))
		if privKeyBlock == nil {
			return "", fmt.Errorf("无法解码 PEM 格式的私钥")
		}
		privateKey, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
		if err != nil {
			// 尝试用 EC 私钥格式解析
			privateKey, err = x509.ParseECPrivateKey(privKeyBlock.Bytes)
			if err != nil {
				return "", fmt.Errorf("解析私钥失败: %w", err)
			}
		}

		// 从私钥派生公钥
		switch key := privateKey.(type) {
		case *ecdsa.PrivateKey:
			publicKey = &key.PublicKey
		default:
			return "", fmt.Errorf("不支持的私钥类型")
		}
	} else {
		return "", fmt.Errorf("必须提供私钥或公钥")
	}

	// 提取未压缩的公钥字节
	switch pub := publicKey.(type) {
	case *ecdsa.PublicKey:
		// 序列化为未压缩格式 (0x04 + X + Y)
		rawPublicKeyBytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		return base64.StdEncoding.EncodeToString(rawPublicKeyBytes), nil
	default:
		return "", fmt.Errorf("不支持的公钥类型")
	}
}
