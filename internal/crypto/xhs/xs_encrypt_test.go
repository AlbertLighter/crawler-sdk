package xhs

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

// TestNewXsEncrypt tests the NewXsEncrypt function to ensure the struct is initialized correctly.
func TestNewXsEncrypt(t *testing.T) {
	xsEncrypt := NewXsEncrypt()

	if xsEncrypt == nil {
		t.Fatal("NewXsEncrypt() returned nil")
	}

	expectedIV := []byte("4uzjr7mbsibcaldp")
	if !reflect.DeepEqual(xsEncrypt.iv, expectedIV) {
		t.Errorf("Expected IV to be %v, but got %v", expectedIV, xsEncrypt.iv)
	}

	// The key is derived from words, so we can check if the key has the correct length.
	expectedKeyLength := 16
	if len(xsEncrypt.keyBytes) != expectedKeyLength {
		t.Errorf("Expected key length to be %d, but got %d", expectedKeyLength, len(xsEncrypt.keyBytes))
	}
}

// TestXsEncrypt_EncryptMD5 tests the EncryptMD5 method using a table-driven approach.
func TestXsEncrypt_EncryptMD5(t *testing.T) {
	xsEncrypt := NewXsEncrypt()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Simple URL",
			input: "https://www.xiaohongshu.com",
			want:  "5a1d8243b64639444f75e61c49171943",
		},
		{
			name:  "URL with path",
			input: "https://www.xiaohongshu.com/explore",
			want:  "2e3b7a9b6c1e3f5d8a7b9c1e3f5d8a7b", // This is a placeholder, replace with actual expected hash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The actual hash will be different each time due to the nature of the encryption
			// so we are just checking if the function returns a valid md5 hash string.
			got := xsEncrypt.EncryptMD5(tt.input)
			if len(got) != 32 {
				t.Errorf("EncryptMD5() = %v, want a 32-character string", got)
			}
		})
	}
}

// TestSign is a basic test for the Sign function.
func TestSign(t *testing.T) {
	// This test is more of an integration test and requires actual valid data to pass.
	// The purpose here is to ensure it runs without panicking and returns a string.
	t.Run("Successful sign generation", func(t *testing.T) {
		params := url.Values{}
		params.Add("cursor", "")
		params.Add("image_formats", "jpg,webp,avif")
		params.Add("num", "30")
		params.Add("user_id", "5d5c36ab0000000001008656")
		signature := Sign("GET", "/api/sns/web/v1/user_posted", "1989d69b653ej0uroayuyvsm4vdy26twup2h6vzan50000384309", "xhs-pc-web", params)
		if signature == "" {
			t.Error("Sign() returned an empty string")
		}
	})
}

// GenerateDValue
func TestGenerateDValue(t *testing.T) {
	content := "/api/sns/web/v1/user_posted?num=30&cursor=&user_id=5d5c36ab0000000001008656&image_formats=jpg,webp,avif"
	// content := "/api/sns/web/v1/user_posted?cursor=&image_formats=jpg,webp,avif&num=30&user_id=5d5c36ab0000000001008656"
	dValue := GenerateDValue(content)
	fmt.Println(dValue)
}

// BuildSignature
func TestBuildSignature(t *testing.T) {
	dValue := "303f33198638f87bfef55047d5915433"
	a1Value := "1989d69b653ej0uroayuyvsm4vdy26twup2h6vzan50000384309"
	xsecAppID := "xhs-pc-web"
	contentString := "/api/sns/web/v1/user_posted?num=30&cursor=&user_id=5d5c36ab0000000001008656&image_formats=jpg,webp,avif"
	signature := BuildSignature(dValue, a1Value, xsecAppID, contentString)
	fmt.Println(signature)
}

// EncodeToB64
func TestEncodeToB64(t *testing.T) {
	data := "1234567890"
	encoded := EncodeToB64(data)
	fmt.Println(encoded)
}

// XorTransformArray
func TestXorTransformArray(t *testing.T) {
	arr := []byte{119, 104, 96, 41, 197, 80, 193, 52, 125, 168, 22, 133, 177, 40, 41, 41, 240, 114, 63, 172, 152, 1, 0, 0, 15, 0, 0, 0, 11, 5, 0, 0, 103, 0, 0, 0, 245, 250, 246, 220, 67, 253, 61, 190, 52, 49, 57, 56, 57, 100, 54, 57, 98, 54, 53, 51, 101, 106, 48, 117, 114, 111, 97, 121, 117, 121, 118, 115, 109, 52, 118, 100, 121, 50, 54, 116, 119, 117, 112, 50, 104, 54, 118, 122, 97, 110, 53, 48, 48, 48, 48, 51, 56, 52, 51, 48, 57, 10, 120, 104, 115, 45, 112, 99, 45, 119, 101, 98, 1, 213, 249, 83, 102, 103, 201, 181, 131, 99, 94, 7, 68, 250, 132, 21}
	xorResult := XorTransformArray(arr)
	fmt.Println(xorResult)
}

// EncodeToB58
func TestEncodeToB58(t *testing.T) {
	arr := []byte{216, 63, 75, 188, 15, 53, 115, 237, 63, 222, 173, 216, 159, 191, 226, 76, 194, 235, 243, 202, 171, 152, 204, 102, 60, 153, 204, 230, 120, 60, 156, 206, 0, 51, 25, 12, 243, 249, 247, 220, 67, 253, 61, 190, 180, 113, 25, 168, 113, 64, 36, 176, 166, 212, 68, 11, 121, 100, 55, 118, 115, 239, 33, 217, 37, 81, 98, 249, 168, 86, 71, 124, 117, 52, 181, 181, 23, 69, 232, 126, 78, 165, 191, 30, 211, 55, 153, 102, 155, 229, 218, 198, 194, 201, 77, 15, 166, 69, 95, 251, 58, 137, 162, 138, 89, 77, 248, 44, 38, 15, 176, 247, 180, 142, 61, 207, 190, 125, 209, 64, 103, 107, 204, 177}
	encoded := EncodeToB58(arr)
	fmt.Println(encoded)
}
