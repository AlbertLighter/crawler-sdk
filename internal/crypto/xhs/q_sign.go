package xhs

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// er:Expires:undefined
// ForceSignHost:true
// Headers:{Content-Length: 16368, Host: 'ros-upload-d4.xhscdn.com'}
// KeyTime:"1755245595;1755331995"
// Method:"PUT"
// Pathname:"/spectrum/7vynuOPF72T6CexEgu_Hcy5c6NcUYaa1DU92in4JeItlMu4"
// Query:{}
// SecretId:"null"
// SecretKey:"null"
// SystemClockOffset:0
// UseRawKey:false
// QSignAuthOptions holds the parameters for generating the Tencent Cloud COS signature.
type QSignAuthOptions struct {
	ForceSignHost     bool
	KeyTime           string
	SecretId          string
	SecretKey         string
	SystemClockOffset int64
	Method            string
	Pathname          string
	Query             map[string]string
	Headers           map[string]string
	Expires           int64 // Expires in seconds
	UseRawKey         bool
}

// getObjectKeys returns the sorted, lowercased keys of a map.
func getObjectKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)
	return keys
}

// obj2str converts a map to a URL-encoded string with sorted keys.
// The keys are lowercased, and the values are URL-encoded.
func obj2str(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		val := m[k]
		pairs = append(pairs, fmt.Sprintf("%s=%s", strings.ToLower(k), url.QueryEscape(val)))
	}
	return strings.Join(pairs, "&")
}

// hmacSha1 computes the HMAC-SHA1 hash of a message with a key and returns it as a hex string.
func hmacSha1(message, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// sha1Hash computes the SHA1 hash of a string and returns it as a hex string.
func sha1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetQSignAuth generates the Tencent Cloud COS authorization string based on the provided options.
func GetQSignAuth(opts QSignAuthOptions) (string, error) {
	if opts.SecretId == "" {
		return "", fmt.Errorf("missing param SecretId")
	}
	if opts.SecretKey == "" {
		return "", fmt.Errorf("missing param SecretKey")
	}

	method := strings.ToLower(opts.Method)
	if method == "" {
		method = "get"
	}

	pathname := opts.Pathname
	if !strings.HasPrefix(pathname, "/") {
		pathname = "/" + pathname
	}

	headers := opts.Headers
	if headers == nil {
		headers = make(map[string]string)
	}
	query := opts.Query
	if query == nil {
		query = make(map[string]string)
	}

	// Time calculation
	startTime := time.Now().Unix()
	expires := opts.Expires
	if expires == 0 {
		expires = 900 // Default 900 seconds
	}
	endTime := startTime + expires

	signTime := fmt.Sprintf("%d;%d", startTime, endTime)
	keyTime := signTime

	// Get lists of header and query param keys
	headerList := getObjectKeys(headers)
	urlParamList := getObjectKeys(query)

	// Signature calculation
	signKey := hmacSha1(keyTime, opts.SecretKey)

	httpString := fmt.Sprintf("%s\n%s\n%s\n%s",
		method,
		pathname,
		obj2str(query),
		obj2str(headers),
	)

	stringToSign := fmt.Sprintf("sha1\n%s\n%s",
		signTime,
		sha1Hash(httpString),
	)

	signature := hmacSha1(stringToSign, signKey)

	// Assemble final authorization string
	authParts := []string{
		"q-sign-algorithm=sha1",
		"q-ak=" + opts.SecretId,
		"q-sign-time=" + signTime,
		"q-key-time=" + keyTime,
		"q-header-list=" + strings.Join(headerList, ";"),
		"q-url-param-list=" + strings.Join(urlParamList, ";"),
		"q-signature=" + signature,
	}

	return strings.Join(authParts, "&"), nil
}
