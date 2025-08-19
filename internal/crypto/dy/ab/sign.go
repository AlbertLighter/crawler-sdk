package ab

import (
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tjfoc/gmsm/sm3"
)

var sObj = map[string]string{
	"s0": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
	"s1": "Dkdpgh4ZKsQB80/Mfvw36XI1R25+WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
	"s2": "Dkdpgh4ZKsQB80/Mfvw36XI1R25-WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
	"s3": "ckdp1h4ZKsUB80/Mfvw36XIgR25+WQAlEi7NLboqYTOPuzmFjJnryx9HVGDaStCe",
	"s4": "Dkdpgh2ZmsQB80/MfvV36XI1R45-WUAlEixNLwoqYTOPuzKFjJnry79HbGcaStCe",
}

// Signer holds configuration for deterministic signing.
type Signer struct {
	fixedTimestamp int64
	randomValues   []float64
	randIdx        int
}

// NewSigner creates a new Signer for deterministic signing.
func NewSigner(fixedTimestamp int64, randomValues []float64) *Signer {
	return &Signer{
		fixedTimestamp: fixedTimestamp,
		randomValues:   randomValues,
	}
}

func (s *Signer) nextRandom() float64 {
	if s.randomValues != nil && len(s.randomValues) > 0 {
		val := s.randomValues[s.randIdx%len(s.randomValues)]
		s.randIdx++
		return val
	}
	return rand.Float64() * 10000
}

func rc4Encrypt(plaintext, key []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(plaintext))
	cipher.XORKeyStream(dst, plaintext)
	return dst, nil
}

// ResultEncrypt 使用自定义字母表对字符串进行 Base64 编码
// 这个实现逻辑上与 Python 版本完全一致
func ResultEncrypt(longStr string, numKey string) (string, error) {
	// 1. 根据 key 获取字母表
	alphabet, ok := sObj[numKey]
	if !ok {
		return "", fmt.Errorf("invalid numKey: %s", numKey)
	}

	// 2. 将输入字符串转换为字节切片
	// Go 的 string 内部是 UTF-8，直接转换为 []byte 可以逐字节处理，
	// 这与 Python 的 .encode('latin-1') 在处理 ASCII 和单字节字符时行为一致。
	longStrBytes := []byte(longStr)
	dataLen := len(longStrBytes)

	// 3. 使用 strings.Builder 高效地构建结果字符串
	// 预分配容量可以进一步提升性能
	var result strings.Builder
	result.Grow((dataLen/3 + 1) * 4)

	// 4. 以 3 字节为步长遍历数据，这是 Base64 的核心逻辑
	for i := 0; i < dataLen; i += 3 {
		// 定义3个字节变量，并处理数据末尾不足3字节的情况
		// 缺失的字节默认为0，这在位运算中是正确的处理方式
		b1, b2, b3 := longStrBytes[i], byte(0), byte(0)

		if i+1 < dataLen {
			b2 = longStrBytes[i+1]
		}
		if i+2 < dataLen {
			b3 = longStrBytes[i+2]
		}

		// 将3个8位字节（24位）合并成一个 uint32 整数
		// 使用 uint32 确保位移操作不会溢出
		longInt := uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3)

		// 将24位整数拆分为4个6位的索引值
		idx1 := (longInt >> 18) & 0x3F
		idx2 := (longInt >> 12) & 0x3F
		idx3 := (longInt >> 6) & 0x3F
		idx4 := longInt & 0x3F

		// 根据索引从字母表中查找字符并写入结果
		result.WriteByte(alphabet[idx1])
		result.WriteByte(alphabet[idx2])

		// 5. 正确处理补位逻辑
		// 如果原始数据块至少有2个字节，才写入第3个编码字符
		if i+1 < dataLen {
			result.WriteByte(alphabet[idx3])
		}
		// 如果原始数据块有3个字节，才写入第4个编码字符
		if i+2 < dataLen {
			result.WriteByte(alphabet[idx4])
		}
	}

	return result.String(), nil
}

func getLongInt(round int, longStr string) int {
	round *= 3
	if round+2 >= len(longStr) {
		return 0
	}
	return (int(longStr[round]) << 16) | (int(longStr[round+1]) << 8) | int(longStr[round+2])
}

func generRandom(random float64, option []int) []byte {
	return []byte{
		byte((int(random) & 255 & 170) | (option[0] & 85)),
		byte((int(random) & 255 & 85) | (option[0] & 170)),
		byte((int(random) >> 8 & 255 & 170) | (option[1] & 85)),
		byte((int(random) >> 8 & 255 & 85) | (option[1] & 170)),
	}
}

func (s *Signer) generateRc4BbStr(urlSearchParams, userAgent, windowEnvStr, suffix string, arguments []int) (string, error) {
	h := sm3.New()
	h.Write([]byte(urlSearchParams + suffix))
	sum1 := h.Sum(nil)
	h.Reset()
	h.Write(sum1)
	urlSearchParamsList := h.Sum(nil)

	h.Reset()
	h.Write([]byte(suffix))
	sum2 := h.Sum(nil)
	h.Reset()
	h.Write(sum2)
	cus := h.Sum(nil)

	key := []byte{byte(0), 1, byte(arguments[2])}
	// key := []byte{byte(0.00390625), 1, byte(arguments[2])}
	rc4Encrypted, err := rc4Encrypt([]byte(userAgent), key)
	if err != nil {
		return "", err
	}
	uaEncrypted, err := ResultEncrypt(string(rc4Encrypted), "s3")
	if err != nil {
		return "", err
	}
	h.Reset()
	h.Write([]byte(uaEncrypted))
	ua := h.Sum(nil)

	startTime := s.fixedTimestamp
	if startTime == 0 {
		startTime = time.Now().UnixMilli()
	}
	endTime := startTime

	b := make(map[int]int64)
	b[8] = 3
	b[10] = endTime
	b[15] = 6383 // aid
	b[16] = startTime
	b[18] = 44

	b[20] = (b[16] >> 24) & 255
	b[21] = (b[16] >> 16) & 255
	b[22] = (b[16] >> 8) & 255
	b[23] = b[16] & 255
	b[24] = int64(uint64(b[16]) / 256 / 256 / 256 / 256)
	b[25] = int64(uint64(b[16]) / 256 / 256 / 256 / 256 / 256)

	b[26] = int64((arguments[0] >> 24) & 255)
	b[27] = int64((arguments[0] >> 16) & 255)
	b[28] = int64((arguments[0] >> 8) & 255)
	b[29] = int64(arguments[0] & 255)

	b[30] = int64((arguments[1] / 256) & 255)
	b[31] = int64((arguments[1] % 256) & 255)
	b[32] = int64((arguments[1] >> 24) & 255)
	b[33] = int64((arguments[1] >> 16) & 255)

	b[34] = int64((arguments[2] >> 24) & 255)
	b[35] = int64((arguments[2] >> 16) & 255)
	b[36] = int64((arguments[2] >> 8) & 255)
	b[37] = int64(arguments[2] & 255)

	b[38] = int64(urlSearchParamsList[21])
	b[39] = int64(urlSearchParamsList[22])

	b[40] = int64(cus[21])
	b[41] = int64(cus[22])

	b[42] = int64(ua[23])
	b[43] = int64(ua[24])

	b[44] = (b[10] >> 24) & 255
	b[45] = (b[10] >> 16) & 255
	b[46] = (b[10] >> 8) & 255
	b[47] = b[10] & 255
	b[48] = b[8]
	b[49] = int64(uint64(b[10]) / 256 / 256 / 256 / 256)
	b[50] = int64(uint64(b[10]) / 256 / 256 / 256 / 256 / 256)

	pageID := 6241
	b[51] = int64(pageID)
	b[52] = int64((pageID >> 24) & 255)
	b[53] = int64((pageID >> 16) & 255)
	b[54] = int64((pageID >> 8) & 255)
	b[55] = int64(pageID & 255)

	aid := 6383
	b[56] = int64(aid)
	b[57] = int64(aid & 255)
	b[58] = int64((aid >> 8) & 255)
	b[59] = int64((aid >> 16) & 255)
	b[60] = int64((aid >> 24) & 255)

	windowEnvList := []byte(windowEnvStr)
	b[64] = int64(len(windowEnvList))
	b[65] = b[64] & 255
	b[66] = (b[64] >> 8) & 255

	b[69] = 0
	b[70] = b[69] & 255
	b[71] = (b[69] >> 8) & 255

	b[72] = b[18] ^ b[20] ^ b[26] ^ b[30] ^ b[38] ^ b[40] ^ b[42] ^ b[21] ^ b[27] ^ b[31] ^ b[35] ^ b[39] ^ b[41] ^ b[43] ^ b[22] ^
		b[28] ^ b[32] ^ b[36] ^ b[23] ^ b[29] ^ b[33] ^ b[37] ^ b[44] ^ b[45] ^ b[46] ^ b[47] ^ b[48] ^ b[49] ^ b[50] ^ b[24] ^
		b[25] ^ b[52] ^ b[53] ^ b[54] ^ b[55] ^ b[57] ^ b[58] ^ b[59] ^ b[60] ^ b[65] ^ b[66] ^ b[70] ^ b[71]

	var bb []byte
	for _, val := range []int{18, 20, 52, 26, 30, 34, 58, 38, 40, 53, 42, 21, 27, 54, 55, 31, 35, 57, 39, 41, 43, 22, 28, 32, 60, 36, 23, 29, 33, 37, 44, 45, 59, 46, 47, 48, 49, 50, 24, 25, 65, 66, 70, 71} {
		bb = append(bb, byte(b[val]))
	}

	// var aa []int64
	// for _, val := range []int{18, 20, 52, 26, 30, 34, 58, 38, 40, 53, 42, 21, 27, 54, 55, 31, 35, 57, 39, 41, 43, 22, 28, 32, 60, 36, 23, 29, 33, 37, 44, 45, 59, 46, 47, 48, 49, 50, 24, 25, 65, 66, 70, 71} {
	// 	aa = append(aa, (b[val]))
	// }
	// fmt.Println(aa)
	// // 创建一个rune切片
	// runes := make([]rune, len(aa))
	// // 将int类型的码点值赋给rune切片
	// for i, code := range aa {
	// 	runes[i] = rune(code)
	// }
	// // 打印16进制字符串
	// hexStr := ""
	// for i := 0; i < len(runes); i++ {
	// 	hexStr += fmt.Sprintf("%02x", runes[i])
	// }
	// fmt.Println(hexStr)
	// fmt.Println(len(runes))

	bb = append(bb, windowEnvList...)
	bb = append(bb, byte(b[72]))

	encryptedBb, err := rc4Encrypt(bb, []byte{121})
	if err != nil {
		return "", err
	}

	return string(encryptedBb), nil
}

func (s *Signer) generateRandomStr() string {
	var randomStrList []byte
	randomStrList = append(randomStrList, generRandom(s.nextRandom(), []int{3, 45})...)
	randomStrList = append(randomStrList, generRandom(s.nextRandom(), []int{1, 0})...)
	randomStrList = append(randomStrList, generRandom(s.nextRandom(), []int{1, 5})...)
	return string(randomStrList)
}

func (s *Signer) Sign(urlSearchParams, userAgent string, arguments []int) (string, error) {
	rc4BbStr, err := s.generateRc4BbStr(
		urlSearchParams,
		userAgent,
		"1536|747|1536|834|0|30|0|0|1536|834|1536|864|1525|747|24|24|Win32",
		"cus",
		arguments,
	)
	if err != nil {
		return "", err
	}
	resultStr := s.generateRandomStr() + rc4BbStr
	// 打印16进制字符串
	hexStr := ""
	for i := 0; i < len(resultStr); i++ {
		hexStr += fmt.Sprintf("%02x", resultStr[i])
	}
	fmt.Println(hexStr)
	encryptedStr, err := ResultEncrypt(resultStr, "s4")
	if err != nil {
		return "", err
	}
	return encryptedStr + "=", nil
}

func (s *Signer) SignDetail(params, userAgent string) (string, error) {
	return s.Sign(params, userAgent, []int{0, 1, 14})
}

func (s *Signer) SignReply(params, userAgent string) (string, error) {
	return s.Sign(params, userAgent, []int{0, 1, 8})
}

// Top-level functions for backward compatibility

func Sign(urlSearchParams, userAgent string, arguments []int) (string, error) {
	s := NewSigner(0, nil)
	return s.Sign(urlSearchParams, userAgent, arguments)
}

func SignDetail(params, userAgent string) (string, error) {
	s := NewSigner(0, nil)
	return s.SignDetail(params, userAgent)
}

func SignReply(params, userAgent string) (string, error) {
	s := NewSigner(0, nil)
	return s.SignReply(params, userAgent)
}

// Sm3Sum remains a standalone function.
func Sm3Sum(data []byte) string {
	h := sm3.New()
	h.Write(data)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
