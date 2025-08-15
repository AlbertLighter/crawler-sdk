package xhs

import (
	"math/big"
	"strings"
)

// Encode 函數將位元組切片編碼為 Base58 字串。
// 這個實現邏輯與您提供的 Python 版本完全對應。
func EncodeToB58(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	// 1. 將位元組轉換為一個大整數
	// 這對應 Python 中的 _bytes_to_number 方法。
	x := new(big.Int).SetBytes(input)

	// 2. 計算前導零的數量
	// 這對應 Python 中的 _count_leading_zeros 方法。
	leadingZeros := 0
	for _, b := range input {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	// 3. 將數字轉換為 Base58 字符
	// 這對應 Python 中的 _number_to_base58_chars 方法。
	base := big.NewInt(58)
	mod := new(big.Int)
	var result []byte

	// 當數字大於 0 時，重複取餘數
	for x.Cmp(big.NewInt(0)) > 0 {
		x.DivMod(x, base, mod) // x = x / 58, mod = x % 58
		result = append(result, Base58Alphabet[mod.Int64()])
	}

	// 4. 添加前導 '1'
	// Base58 中的 '1' 代表位元組中的前導零。
	for i := 0; i < leadingZeros; i++ {
		result = append(result, Base58Alphabet[0])
	}

	// 5. 反轉結果
	// 因為我們是從最低位開始計算的，所以需要反轉順序。
	reverse(result)

	return string(result)
}

// reverse 函數用於原地反轉位元組切片。
func reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
// EncodeToB64 encodes a string to a Base64 string using the custom alphabet.
func EncodeToB64(data string) string {
	var result strings.Builder
	for _, char := range data {
		if idx := strings.IndexRune(StandardBase64Alphabet, char); idx != -1 {
			result.WriteByte(CustomBase64Alphabet[idx])
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}
