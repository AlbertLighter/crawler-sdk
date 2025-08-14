package xhs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// XscEncrypt 提供了字符串加密与Base64编码的功能
type XscEncrypt struct{}

// EncryptEncodeUTF8 对输入的文本进行URL编码，转换百分号编码为十进制ASCII值
func (x *XscEncrypt) EncryptEncodeUTF8(text string) ([]byte, error) {
	encoded := url.QueryEscape(text)
	var result []byte
	for i := 0; i < len(encoded); {
		if encoded[i] == '%' && i+2 < len(encoded) {
			val, err := strconv.ParseInt(encoded[i+1:i+3], 16, 8)
			if err != nil {
				return nil, fmt.Errorf("解析百分号编码失败: %w", err)
			}
			result = append(result, byte(val))
			i += 3
		} else {
			result = append(result, encoded[i])
			i++
		}
	}
	return result, nil
}

// TripletToBase64 将24位整数分成4个6位部分，转换为Base64字符串
func (x *XscEncrypt) TripletToBase64(e int) string {
	return string(Lookup[(e>>18)&63]) + string(Lookup[(e>>12)&63]) +
		string(Lookup[(e>>6)&63]) + string(Lookup[e&63])
}

// EncodeChunk 将编码后的整数列表分成3字节一组转换为Base64
func (x *XscEncrypt) EncodeChunk(e []byte, t, r int) string {
	var chunks strings.Builder
	for b := t; b < r; b += 3 {
		if b+2 < len(e) { // 确保有完整的三个字节
			chunk := x.TripletToBase64((int(e[b]) << 16) + (int(e[b+1]) << 8) + int(e[b+2]))
			chunks.WriteString(chunk)
		}
	}
	return chunks.String()
}

// B64Encode 将字节切片编码为Base64格式
func (x *XscEncrypt) B64Encode(e []byte) string {
	P := len(e)
	W := P % 3
	Z := P - W
	var result strings.Builder

	for i := 0; i < Z; i += 16383 {
		end := i + 16383
		if end > Z {
			end = Z
		}
		result.WriteString(x.EncodeChunk(e, i, end))
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

// Mrc 使用自定义CRC算法生成校验值
func (x *XscEncrypt) Mrc(e string) uint32 {
	o := uint32(0xFFFFFFFF) // -1 in signed 32-bit is 0xFFFFFFFF unsigned

	for _, char := range e {
		o = IE[(o&0xFF)^uint32(char)] ^ (o >> 8)
	}
	return ^o ^ 3988292384
}

// EncryptXsc 生成xsc
func (x *XscEncrypt) EncryptXsc(xs, xt, platform, a1, x1, x4, b1 string) ([]byte, error) {
	x9 := strconv.FormatUint(uint64(x.Mrc(xt+xs+b1)), 10)

	// x10 is hardcoded to 24 in the Python example, so we'll use that.
	// If it needs to be random, uncomment the rand.Intn line.
	// x10 := rand.Intn(20) + 10 // random.randint(10, 29)
	x10 := 24

	obj := map[string]interface{}{
		"s0":  5,
		"s1":  "",
		"x0":  "1",
		"x1":  x1,
		"x2":  "Windows",
		"x3":  platform,
		"x4":  x4,
		"x5":  a1,
		"x6":  xt,
		"x7":  xs,
		"x8":  b1,
		"x9":  x9,
		"x10": x10,
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %w", err)
	}

	return x.EncryptEncodeUTF8(string(jsonBytes))
}
