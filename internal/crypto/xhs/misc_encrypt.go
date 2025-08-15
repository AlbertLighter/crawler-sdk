package xhs

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// XorTransformArray performs an XOR transformation on a byte array.
func XorTransformArray(source []byte) []byte {
	result := make([]byte, len(source))
	for i := range source {
		result[i] = source[i] ^ HexKey[i]
	}
	return result
}

// IntToLeBytes converts an integer to a little-endian byte array.
func IntToLeBytes(val, length int) []byte {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(val & 0xFF)
		val >>= 8
	}
	return bytes
}

// StrToLenPrefixedBytes converts a string to a byte array with a 1-byte length prefix.
func StrToLenPrefixedBytes(s string) []byte {
	buf := []byte(s)
	return append([]byte{byte(len(buf))}, buf...)
}

// HexStringToBytes converts a hex string to a byte array.
func HexStringToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// ProcessHexParameter processes a hex parameter string.
func ProcessHexParameter(hexStr string, xorKey byte) ([]byte, error) {
	if len(hexStr) != ExpectedHexLength {
		return nil, nil // Return nil to indicate an error in a more Go-idiomatic way.
	}

	byteValues, err := HexStringToBytes(hexStr)
	if err != nil {
		return nil, err
	}

	result := make([]byte, OutputByteCount)
	for i := 0; i < OutputByteCount; i++ {
		result[i] = byteValues[i] ^ xorKey
	}
	return result, nil
}

func Uint32ToLeBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	return b
}

func Uint64ToLeBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

// CustomFieldDecrypt 提供了自定义的字段解密相关功能
type CustomFieldDecrypt struct{}

// RandomStr 生成指定长度的随机字母数字字符串
func (c *CustomFieldDecrypt) RandomStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// Base36Encode 将数字转换为Base36编码
func (c *CustomFieldDecrypt) Base36Encode(number int64) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if number == 0 {
		return string(alphabet[0])
	}

	sign := ""
	if number < 0 {
		sign = "-"
		number = -number
	}

	base36 := ""
	for number > 0 {
		i := number % int64(len(alphabet))
		base36 = string(alphabet[i]) + base36
		number /= int64(len(alphabet))
	}
	return sign + base36
}

// TripletToBase64 将24位整数分成4个6位部分，转换为Base64字符串
func (c *CustomFieldDecrypt) TripletToBase64(e int) string {
	return string(Lookup[(e>>18)&63]) + string(Lookup[(e>>12)&63]) +
		string(Lookup[(e>>6)&63]) + string(Lookup[e&63])
}

// EncodeChunk 将字节切片分成3字节一组转换为Base64
func (c *CustomFieldDecrypt) EncodeChunk(e []byte, t, r int) string {
	var m strings.Builder
	for b := t; b < r; b += 3 {
		if b+2 < len(e) { // Ensure there are full three bytes
			n := (int(e[b]) << 16) + (int(e[b+1]) << 8) + int(e[b+2])
			m.WriteString(c.TripletToBase64(n))
		}
	}
	return m.String()
}

// B64Encode 将字节切片编码为自定义Base64格式
func (c *CustomFieldDecrypt) B64Encode(e []byte) string {
	P := len(e)
	W := P % 3
	Z := P - W
	var result strings.Builder

	for i := 0; i < Z; i += 16383 {
		end := i + 16383
		if end > Z {
			end = Z
		}
		result.WriteString(c.EncodeChunk(e, i, end))
	}

	if W == 1 {
		F := e[P-1]
		result.WriteByte(Lookup[F>>2])
		result.WriteByte(Lookup[(F<<4)&63])
		result.WriteString("==")
	} else if W == 2 {
		F := (int(e[P-2]) << 8) + int(e[P-1])
		result.WriteByte(Lookup[F>>10])
		result.WriteByte(Lookup[(F>>4)&63])
		result.WriteByte(Lookup[(F<<2)&63])
		result.WriteString("=")
	}
	return result.String()
}

// CookieFieldEncrypt 提供了Cookie字段加密功能
type CookieFieldEncrypt struct {
	cfd CustomFieldDecrypt
}

// GetA1AndWebID 生成 a1 和 webid
func (c *CookieFieldEncrypt) GetA1AndWebID() (string, string) {
	d := fmt.Sprintf("%x%s50000", time.Now().UnixNano()/1e6, c.cfd.RandomStr(30))
	crc := crc32.ChecksumIEEE([]byte(d))
	g := d + strconv.FormatUint(uint64(crc), 10)
	if len(g) > 52 {
		g = g[:52]
	}

	hasher := md5.New()
	hasher.Write([]byte(g))
	webID := hex.EncodeToString(hasher.Sum(nil))

	return g, webID
}

// MiscEncrypt 提供了其他杂项加密功能
type MiscEncrypt struct {
	cfd CustomFieldDecrypt
	cfe CookieFieldEncrypt
}

// NewMiscEncrypt 创建 MiscEncrypt 实例
func NewMiscEncrypt() *MiscEncrypt {
	cfd := CustomFieldDecrypt{}
	return &MiscEncrypt{
		cfd: cfd,
		cfe: CookieFieldEncrypt{cfd: cfd},
	}
}

// X_B3_TraceID 生成 x_b3_traceid
func (m *MiscEncrypt) X_B3_TraceID() string {
	const characters = "abcdef0123456789"
	sb := strings.Builder{}
	sb.Grow(16)
	for i := 0; i < 16; i++ {
		sb.WriteByte(characters[rand.Intn(len(characters))])
	}
	return sb.String()
}

// SearchID 生成 search_id
func (m *MiscEncrypt) SearchID() string {
	e := time.Now().UnixNano() / 1e6 << 64 // This shift might be problematic for int64, check Python's behavior
	t := rand.Int63n(2147483646)           // Max int for Python's random.uniform(0, 2147483646)
	return m.cfd.Base36Encode(e + t)
}

// X_Xray_TraceID 生成 x_xray_traceid
func (m *MiscEncrypt) X_Xray_TraceID(x_b3 string) string {
	hasher := md5.New()
	hasher.Write([]byte(x_b3))
	return hex.EncodeToString(hasher.Sum(nil))
}
