package ab

import (
	"testing"
)

func TestSignConsistency(t *testing.T) {
	// These parameters are the same as in the original JS test environment.
	urlSearchParams := "device_platform=webapp&aid=6383&channel=channel_pc_web&update_version_code=170400&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1536&screen_height=864&browser_language=zh-CN&browser_platform=Win32&browser_name=Chrome&browser_version=123.0.0.0&browser_online=true&engine_name=Blink&engine_version=123.0.0.0&os_name=Windows&os_version=10&cpu_core_num=16&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=50&webid=7362810250930783783&msToken=VkDUvz1y24CppXSl80iFPr6ez-3FiizcwD7fI1OqBt6IICq9RWG7nCvxKb8IVi55mFd-wnqoNkXGnxHrikQb4PuKob5Q-YhDp5Um215JzlBszkUyiEvR"
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"

	// The expected signature is obtained by running the js_helper.js script with a fixed seed.
	expectedSignature := "xyRhBmhfDk2p6DS65I2LfY3q6fN3YgbA0trEMD2fpVVWiL39HMYD9exoWN4v3Y8joT/IIeYjy4hbT3ohrQ2y8qwf9W0L/25gsDSkKl12so0j53inCLf/E0iE5hsAtFH8svr4iKi8owICSYyhldAJ5kIlO62-zo0/9-j="

	// Use the same fixed values as in the JS helper
	fixedTimestamp := int64(1678886400000)
	// The JS code uses Math.random() * 10000, so we provide the scaled values.
	randomValues := []float64{1230, 4560, 7890}

	// Create a new signer with deterministic values.
	signer := NewSigner(fixedTimestamp, randomValues)

	actualSignature, err := signer.SignDetail(urlSearchParams, userAgent)
	if err != nil {
		t.Fatalf("SignDetail returned an error: %v", err)
	}
	for i := 0; i < len(actualSignature); i++ {
		if actualSignature[i] != expectedSignature[i] {
			t.Errorf("Signature mismatch at index %d:\n", i)
		}
	}
	if actualSignature != expectedSignature {
		t.Errorf("Signature mismatch:\nExpected: %s\nActual:   %s", expectedSignature, actualSignature)
	}
}