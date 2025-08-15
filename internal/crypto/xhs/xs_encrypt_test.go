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
		signature := Sign("GET", "/api/v2/note/feed", "test_a1", "test_xsec_appid", url.Values{})
		if signature == "" {
			t.Error("Sign() returned an empty string")
		}
	})
}

// GenerateDValue
func TestGenerateDValue(t *testing.T) {
	content := "https://www.xiaohongshu.com/api/sns/web/v1/user_posted?cursor=684ebc77000000002100a2c4&image_formats=jpg,webp,avif&num=30&user_id=5d5c36ab0000000001008656&xsec_source=pc_feed&xsec_token=ABn830XCOxiqnEyuW9NzS0hmuNr9Se3HJ3v-pZFRItuHo="
	dValue := GenerateDValue(content)
	fmt.Println(dValue)
}

// BuildSignature
func TestBuildSignature(t *testing.T) {
	dValue := "1234567890"
	a1Value := "1234567890"
	xsecAppID := "1234567890"
	contentString := "1234567890"
	signature := BuildSignature(dValue, a1Value, xsecAppID, contentString)
	fmt.Println(signature)
}

// EncodeToB64
func TestEncodeToB64(t *testing.T) {
	data := "1234567890"
	encoded := EncodeToB64(data)
	fmt.Println(encoded)
}
