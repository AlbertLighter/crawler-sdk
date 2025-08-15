package xhs

import (
	"math/big"
	"strings"
)

// EncodeToB58 encodes a byte array to a Base58 string using the custom alphabet.
func EncodeToB58(input []byte) string {
	num := new(big.Int).SetBytes(input)
	base := big.NewInt(int64(len(Base58Alphabet)))
	mod := &big.Int{}

	var result []byte
	for num.Cmp(big.NewInt(0)) > 0 {
		num.DivMod(num, base, mod)
		result = append(result, Base58Alphabet[mod.Int64()])
	}

	// Add leading zeros
	for _, b := range input {
		if b == 0 {
			result = append(result, Base58Alphabet[0])
		} else {
			break
		}
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
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
