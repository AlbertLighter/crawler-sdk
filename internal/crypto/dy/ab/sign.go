package ab

import (
	"crypto/rc4"
	"encoding/hex"
	"math"
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

func resultEncrypt(longStr string, num string) string {
	sMap := map[string]interface{}{
		"0":   16515072,
		"1":   258048,
		"2":   4032,
		"str": sObj[num],
	}

	var result strings.Builder
	lound := 0
	longInt := getLongInt(lound, longStr)

	for i := 0; i < len(longStr)/3*4; i++ {
		if math.Floor(float64(i)/4) != float64(lound) {
			lound++
			longInt = getLongInt(lound, longStr)
		}
		key := i % 4
		var tempInt int
		switch key {
		case 0:
			tempInt = (longInt & sMap["0"].(int)) >> 18
			result.WriteByte(sMap["str"].(string)[tempInt])
		case 1:
			tempInt = (longInt & sMap["1"].(int)) >> 12
			result.WriteByte(sMap["str"].(string)[tempInt])
		case 2:
			tempInt = (longInt & sMap["2"].(int)) >> 6
			result.WriteByte(sMap["str"].(string)[tempInt])
		case 3:
			tempInt = longInt & 63
			result.WriteByte(sMap["str"].(string)[tempInt])
		}
	}
	return result.String()
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

	key := []byte{byte(0.00390625 * 256), 1, byte(arguments[2])}
	rc4Encrypted, err := rc4Encrypt([]byte(userAgent), key)
	if err != nil {
		return "", err
	}
	uaEncrypted := resultEncrypt(string(rc4Encrypted), "s3")
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
	return resultEncrypt(resultStr, "s4") + "=", nil
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