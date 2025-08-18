package dy

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"
)

func hexToBigInt(hexStr string) *big.Int {
	n := new(big.Int)
	n.SetString(hexStr, 16)
	return n
}

func genEcdsaFromHex(hexString string) *KeyManage {
	manage := new(KeyManage)

	pubKeyHex := hexString[:65*2]

	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		log.Println(err)
		return nil
	}

	manage.pubKeyBaseString = base64.StdEncoding.EncodeToString(pubKeyBytes)

	priKeyHex := hexString[65*2:]

	privateKeyInt := hexToBigInt(priKeyHex)

	// 创建ECDSA私钥对象
	curve := elliptic.P256()
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = curve
	privateKey.D = privateKeyInt
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyInt.Bytes())

	//fmt.Printf("Private Key: %x\n", privateKey.D)
	//fmt.Printf("Public Key: %x\n", elliptic.Marshal(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y))

	manage.privateKey = privateKey

	return manage
}

func genEcdsa() (manage *KeyManage, hexString string) {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println(err)
		return
	}
	//打印私钥
	//fmt.Printf("Private Key len = %d: %x\n", len(priKey.D.Bytes()), priKey.D.Bytes())

	hexString = priKey.D.Text(16)

	pubKeyBytes := elliptic.Marshal(elliptic.P256(), priKey.X, priKey.Y)

	//fmt.Println(hex.EncodeToString(pubKeyBytes))

	hexString = hex.EncodeToString(pubKeyBytes) + hexString

	manage = new(KeyManage)
	manage.privateKey = priKey

	manage.pubKeyBaseString = base64.StdEncoding.EncodeToString(pubKeyBytes)

	return
}

func genCsr(priKey *ecdsa.PrivateKey) string {
	subj := pkix.Name{
		CommonName: "bd-ticket-guard",
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrBytes, e := x509.CreateCertificateRequest(rand.Reader, &template, priKey)
	if e != nil {
		log.Println(e)
	}

	block := &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}

	buf := new(bytes.Buffer)
	pem.Encode(buf, block)
	//pem.Encode(os.Stdout, &pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})

	m := pem.EncodeToMemory(block) // []byte

	return base64.StdEncoding.EncodeToString(m)
}

func genCert(priKey *ecdsa.PrivateKey) string {
	subj := pkix.Name{
		CommonName: "bd-ticket-guard",
	}

	// 生成一个随机序列号
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber:       serialNumber,
		Subject:            subj,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().AddDate(1, 0, 0), // 证书有效期为1年
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priKey.PublicKey, priKey)
	if err != nil {
		log.Println(err)
	}

	block := &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	buf := new(bytes.Buffer)
	pem.Encode(buf, block)

	encoded := pem.EncodeToMemory(block) // []byte

	return base64.StdEncoding.EncodeToString(encoded)
}

type KeyManage struct {
	pubKeyBaseString string
	privateKey       *ecdsa.PrivateKey
}

func (key *KeyManage) Sign(content string) string {
	if key.privateKey == nil {
		return ""
	}

	hash := sha256.Sum256([]byte(content))

	sig, e := ecdsa.SignASN1(rand.Reader, key.privateKey, hash[:])
	if e != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(sig)
}

type BDTGKeyPair struct {
	TsSignTee        string `json:"ts_sign_tee"`
	TsSignRee        string `json:"ts_sign_ree"`
	TeePrivateKeyHex string `json:"tee_private_key"`
	ReePrivateKeyHex string `json:"ree_private_key"`
	XTTToken         string `json:"x_tt_token"`
	ClientCert       string `json:"client_cert"`
	ServerCert       string `json:"server_cert"`
	ServerSn         string `json:"server_sn"`

	teeManage *KeyManage
	reeManage *KeyManage
}

func (bd *BDTGKeyPair) GetReePubKeyBase64String() string {
	if bd.reeManage != nil {
		return bd.reeManage.pubKeyBaseString
	}
	return ""
}

func (bd *BDTGKeyPair) BuildReeKey() bool {
	if len(bd.ReePrivateKeyHex) != 97*2 {
		return false
	}

	bd.reeManage = genEcdsaFromHex(bd.ReePrivateKeyHex)

	if bd.reeManage != nil {
		return true
	} else {
		return false
	}
}

func (bd *BDTGKeyPair) BuildTeeKey() bool {
	if len(bd.TeePrivateKeyHex) == 97*2 {
		bd.teeManage = genEcdsaFromHex(bd.TeePrivateKeyHex)
	} else {
		bd.teeManage, bd.TeePrivateKeyHex = genEcdsa()
	}

	if bd.teeManage != nil {
		return true
	} else {
		return false
	}
}

func (bd *BDTGKeyPair) GetClientCsr() string {
	if bd.teeManage == nil || bd.teeManage.privateKey == nil {
		return ""
	}
	return genCsr(bd.teeManage.privateKey)
}

func (bd *BDTGKeyPair) GetClientCert() string {
	if bd.teeManage == nil || bd.teeManage.privateKey == nil {
		return ""
	}
	return genCert(bd.teeManage.privateKey)
}

func (bd *BDTGKeyPair) GetPostCsr() string {
	if bd.reeManage == nil || bd.teeManage.privateKey == nil {
		return ""
	}
	return genCsr(bd.reeManage.privateKey)
}

func (bd *BDTGKeyPair) GetClientData(path string) string {
	if bd.reeManage == nil ||
		bd.reeManage.privateKey == nil ||
		bd.teeManage == nil ||
		bd.teeManage.privateKey == nil {
		//return ""
	}
	if len(bd.TsSignRee) < 10 || len(bd.TsSignTee) < 10 {
		return ""
	}

	clientInfo := make(map[string]interface{})
	clientInfo["req_content"] = "ticket,path,timestamp"
	//clientInfo["ts_sign"] = bd.TsSignTee
	clientInfo["ts_sign_ree"] = bd.TsSignRee

	now := time.Now().Unix()

	info := fmt.Sprintf("ticket=%s&path=%s&timestamp=%d", bd.XTTToken, path, now)

	clientInfo["timestamp"] = now
	//clientInfo["req_sign"] = bd.teeManage.Sign(info)
	clientInfo["req_sign_ree"] = bd.reeManage.Sign(info)

	clientData, err := json.Marshal(clientInfo)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(clientData)
}

// func Config(b *DouyinApi) flurl.Config {
// 	return func(rb *flurl.Builder) {
// 		u, err := rb.URL()
// 		if err != nil {
// 			rb.SetErr(err)
// 			return
// 		}
// 		urlPath := u.Path

// 		k := &BDTGKeyPair{
// 			TsSignRee:        b.TsSignRee,
// 			TsSignTee:        b.TsSignTee,
// 			ReePrivateKeyHex: b.ReePrivateKeyHex,
// 			XTTToken:         b.Token,
// 		}

// 		// 解析URL
// 		n, _ := strconv.Atoi(b.BuildNumber)
// 		if n < 298014 {
// 			//logger.Printf("app_version < 29 账号：%s", b.UserName)
// 			return
// 		} else if k.ReePrivateKeyHex == "" || k.TsSignRee == "" || k.TsSignTee == "" {
// 			log.Println("BDSign Key Error NULL 账号：", b.UserName)
// 			//rb.SetErr(fmt.Errorf("BDSign Key Error 账号：%s", b.UserName))
// 			return
// 		}
// 		if !k.BuildTeeKey() || !k.BuildReeKey() {
// 			rb.SetErr(fmt.Errorf("BDSign Build Key Error 账号：%s", b.UserName))
// 			return
// 		}
// 		headers := map[string]string{
// 			"bd-ticket-guard-client-data":       k.GetClientData(urlPath),
// 			"bd-ticket-guard-ree-public-key":    k.GetReePubKeyBase64String(),
// 			"bd-ticket-guard-client-cert":       k.GetClientCert(),
// 			"bd-ticket-guard-version":           "3",
// 			"bd-ticket-guard-iteration-version": "2",
// 		}
// 		// 将headers添加到请求中
// 		for key, value := range headers {
// 			if value == "" {
// 				rb.SetErr(fmt.Errorf("BDSign Key Error 账号：%s key:%s", b.UserName, key))
// 				return
// 			}
// 			rb.Header(key, value)
// 		}
// 	}
// }

// func GetBdSign(b *DouyinApi, u string) map[string]string {
// 	k := &BDTGKeyPair{
// 		TsSignRee:        b.TsSignRee,
// 		TsSignTee:        b.TsSignTee,
// 		ReePrivateKeyHex: b.ReePrivateKeyHex,
// 		XTTToken:         b.Token,
// 	}

// 	// 解析URL
// 	n, _ := strconv.Atoi(b.BuildNumber)
// 	if n < 298014 {
// 		//logger.Printf("app_version < 29 账号：%s", b.UserName)
// 		return nil
// 	} else if k.ReePrivateKeyHex == "" || k.TsSignRee == "" || k.TsSignTee == "" {
// 		log.Println("BDSign Key Error NULL 账号：", b.UserName)
// 		//rb.SetErr(fmt.Errorf("BDSign Key Error 账号：%s", b.UserName))
// 		return nil
// 	}
// 	if !k.BuildTeeKey() || !k.BuildReeKey() {
// 		return nil
// 	}
// 	headers := map[string]string{
// 		"bd-ticket-guard-client-data":       k.GetClientData(u),
// 		"bd-ticket-guard-ree-public-key":    k.GetReePubKeyBase64String(),
// 		"bd-ticket-guard-client-cert":       k.GetClientCert(),
// 		"bd-ticket-guard-version":           "3",
// 		"bd-ticket-guard-iteration-version": "2",
// 	}
// 	return headers
// }
