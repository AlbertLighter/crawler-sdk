package dy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Constants for AWS v4 Signing
const (
	awsAlgorithm           = "AWS4-HMAC-SHA256"
	awsV4Identifier        = "aws4_request"
	awsDateHeader          = "X-Amz-Date"
	awsTokenHeader         = "X-Amz-Security-Token"
	awsContentSha256Header = "X-Amz-Content-Sha256"
	awsKDatePrefix         = "AWS4"
)

// Signer holds the credentials and configuration for signing requests.
// It's designed to be compatible with AWS Signature Version 4.
type Signer struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	Service         string
}

// NewSigner creates a new Signer.
func NewSigner(accessKey, secretKey, sessionToken, region, service string) *Signer {
	return &Signer{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		SessionToken:    sessionToken,
		Region:          region,
		Service:         service,
	}
}

// Sign calculates and adds the AWS v4 authorization headers to the given http.Request.
func (s *Signer) Sign(req *http.Request, body []byte, signTime time.Time) {
	// 1. Prepare timestamps
	amzDate := signTime.UTC().Format("20060102T150405Z")

	dateStamp := signTime.UTC().Format("20060102")

	// 2. Calculate payload hash

	payloadHash := hex.EncodeToString(hashSHA256(body))
	req.Header.Set(awsContentSha256Header, payloadHash)
	req.Header.Set(awsDateHeader, amzDate)
	if s.SessionToken != "" {
		req.Header.Set(awsTokenHeader, s.SessionToken)
	}

	// 3. Create Canonical Request

	canonicalRequest, signedHeaders := s.createCanonicalRequest(req, payloadHash)

	// 4. Create String to Sign

	credentialScope := fmt.Sprintf("%s/%s/%s/%s", dateStamp, s.Region, s.Service, awsV4Identifier)
	stringToSign := s.createStringToSign(amzDate, credentialScope, canonicalRequest)

	// 5. Calculate Signature

	signingKey := s.getSignatureKey(dateStamp)
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// 6. Build Authorization Header

	authHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		awsAlgorithm,
		s.AccessKeyID,
		credentialScope,
		signedHeaders,
		signature)
	req.Header.Set("Authorization", authHeader)
}

// createCanonicalRequest creates the canonical request string.
func (s *Signer) createCanonicalRequest(req *http.Request, payloadHash string) (string, string) {
	// Canonical URI
	canonicalURI := req.URL.Path
	if canonicalURI == "" {
		canonicalURI = "/"
	}

	// Canonical Query String

	canonicalQuery := s.createCanonicalQuery(req.URL.Query())

	// Canonical Headers & Signed Headers

	canonicalHeaders, signedHeaders := s.createCanonicalHeaders(req)

	return fmt.Sprintf("%s\n%s\n%s\n%s\n\n%s\n%s",
		req.Method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	), signedHeaders
}

// createCanonicalQuery sorts and encodes the query parameters.
func (s *Signer) createCanonicalQuery(queryParams url.Values) string {
	keys := make([]string, 0, len(queryParams))
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, key := range keys {
		values := queryParams[key]
		sort.Strings(values)
		for _, value := range values {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}
	}
	return strings.Join(queryParts, "&")
}

// createCanonicalHeaders creates the canonical headers string and the signed headers list.
func (s *Signer) createCanonicalHeaders(req *http.Request) (string, string) {
	var headers [][2]string
	for key, values := range req.Header {
		lowerKey := strings.ToLower(key)
		// All headers used for signing must be lowercased.
		headers = append(headers, [2]string{lowerKey, strings.TrimSpace(strings.Join(values, ","))})
	}
	// Sort by header name
	sort.Slice(headers, func(i, j int) bool {
		return headers[i][0] < headers[j][0]
	})

	var canonicalHeadersParts []string
	var signedHeaderParts []string
	for _, header := range headers {
		canonicalHeadersParts = append(canonicalHeadersParts, header[0]+":"+header[1])
		signedHeaderParts = append(signedHeaderParts, header[0])
	}

	return strings.Join(canonicalHeadersParts, "\n"), strings.Join(signedHeaderParts, ";")
}

// createStringToSign creates the string that will be signed.
func (s *Signer) createStringToSign(amzDate, credentialScope, canonicalRequest string) string {
	canonicalRequestHash := hex.EncodeToString(hashSHA256([]byte(canonicalRequest)))
	return fmt.Sprintf("%s\n%s\n%s\n%s",
		awsAlgorithm,
		amzDate,
		credentialScope,
		canonicalRequestHash,
	)
}

// getSignatureKey calculates the derived signing key.
func (s *Signer) getSignatureKey(dateStamp string) []byte {
	kDate := hmacSHA256([]byte(awsKDatePrefix+s.SecretAccessKey), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(s.Region))
	kService := hmacSHA256(kRegion, []byte(s.Service))
	kSigning := hmacSHA256(kService, []byte(awsV4Identifier))
	return kSigning
}

// Helper functions for hashing
func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func hmacSHA256(key []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
