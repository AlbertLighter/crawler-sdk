package dy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"resty.dev/v3"
)

// // Constants for AWS v4 Signing
// const (
// 	awsAlgorithm           = "AWS4-HMAC-SHA256"
// 	awsV4Identifier        = "aws4_request"
// 	awsDateHeader          = "X-Amz-Date"
// 	awsTokenHeader         = "X-Amz-Security-Token"
// 	awsContentSha256Header = "X-Amz-Content-Sha256"
// 	awsKDatePrefix         = "AWS4"
// )

// // Signer holds the credentials and configuration for signing requests.
// // It's designed to be compatible with AWS Signature Version 4.
// type Signer struct {
// 	AccessKeyID     string
// 	SecretAccessKey string
// 	SessionToken    string
// 	Region          string
// 	Service         string
// }

// // NewSigner creates a new Signer.
// func NewSigner(accessKey, secretKey, sessionToken, region, service string) *Signer {
// 	return &Signer{
// 		AccessKeyID:     accessKey,
// 		SecretAccessKey: secretKey,
// 		SessionToken:    sessionToken,
// 		Region:          region,  //cn-north-1
// 		Service:         service, //imagex
// 	}
// }

// Signer holds the signing information.
type Signer struct {
	Request             *resty.Request
	ServiceName         string //imagex
	Region              string //cn-north-1
	Constant            signerConst
	BodySha256          string
	ShouldSerializeBody bool
}

type signerConst struct {
	Algorithm           string
	V4Identifier        string
	DateHeader          string
	TokenHeader         string
	ContentSha256Header string
	KDatePrefix         string
}

// NewSigner creates a new Signer.
func NewSigner(req *http.Request, serviceName, region string, isVolcengine bool) *Signer {
	s := &Signer{
		Request:             req,
		ServiceName:         serviceName,
		Region:              region,
		ShouldSerializeBody: true,
	}
	if isVolcengine {
		s.Constant = signerConst{
			Algorithm:           "HMAC-SHA256",
			V4Identifier:        "request",
			DateHeader:          "X-Date",
			TokenHeader:         "x-security-token",
			ContentSha256Header: "X-Content-Sha256",
			KDatePrefix:         "",
		}
	} else {
		s.Constant = signerConst{
			Algorithm:           "AWS4-HMAC-SHA256",
			V4Identifier:        "aws4_request",
			DateHeader:          "X-Amz-Date",
			TokenHeader:         "x-amz-security-token",
			ContentSha256Header: "X-Amz-Content-Sha256",
			KDatePrefix:         "AWS4",
		}
	}
	return s
}

// AddAuthorization adds the Authorization header to the request.
func (s *Signer) AddAuthorization(credentials map[string]string, date time.Time) {
	isoDate := s.iso8601(date)
	s.addHeaders(credentials, isoDate)
	authHeader := s.authorization(credentials, isoDate)
	s.Request.Header.Set("Authorization", authHeader)
}

func (s *Signer) addHeaders(credentials map[string]string, isoDate string) {
	s.Request.Header.Set(s.Constant.DateHeader, isoDate)
	if sessionToken, ok := credentials["sessionToken"]; ok {
		s.Request.Header.Set(s.Constant.TokenHeader, sessionToken)
	}
	if s.Request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(s.Request.Body)
		s.Request.Header.Set(s.Constant.ContentSha256Header, s.hexEncodedHash(string(bodyBytes)))
	}
}

func (s *Signer) authorization(credentials map[string]string, isoDate string) string {
	credentialString := s.credentialString(isoDate)
	signedHeaders := s.signedHeaders()
	signature := s.signature(credentials, isoDate)

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.Constant.Algorithm,
		credentials["accessKeyId"],
		credentialString,
		signedHeaders,
		signature)
}

func (s *Signer) signature(credentials map[string]string, isoDate string) string {
	signingKey := s.getSigningKey(credentials, isoDate[0:8], s.Region, s.ServiceName)
	stringToSign := s.stringToSign(isoDate)

	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(stringToSign))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Signer) stringToSign(isoDate string) string {
	return strings.Join([]string{
		s.Constant.Algorithm,
		isoDate,
		s.credentialString(isoDate),
		s.hexEncodedHash(s.canonicalString()),
	}, "\n")
}

func (s *Signer) canonicalString() string {
	return strings.Join([]string{
		s.Request.Method,
		s.Request.URL.Path,
		s.canonicalQueryString(),
		s.canonicalHeaders(),
		s.signedHeaders(),
		s.hexEncodedBodyHash(),
	}, "\n")
}

func (s *Signer) canonicalQueryString() string {
	var keys []string
	query := s.Request.URL.Query()
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var params []string
	for _, k := range keys {
		values := query[k]
		sort.Strings(values)
		for _, v := range values {
			params = append(params, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}
	return strings.Join(params, "&")
}

func (s *Signer) canonicalHeaders() string {
	var headers []string
	var headerKeys []string
	for k := range s.Request.Header {
		headerKeys = append(headerKeys, strings.ToLower(k))
	}
	sort.Strings(headerKeys)

	for _, k := range headerKeys {
		if s.isSignableHeader(k) {
			headers = append(headers, k+":"+strings.TrimSpace(s.Request.Header.Get(k)))
		}
	}
	return strings.Join(headers, "\n")
}

func (s *Signer) signedHeaders() string {
	var headerKeys []string
	for k := range s.Request.Header {
		if s.isSignableHeader(strings.ToLower(k)) {
			headerKeys = append(headerKeys, strings.ToLower(k))
		}
	}
	sort.Strings(headerKeys)
	return strings.Join(headerKeys, ";")
}

func (s *Signer) credentialString(isoDate string) string {
	return s.createScope(isoDate[0:8], s.Region, s.ServiceName)
}

func (s *Signer) hexEncodedHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *Signer) hexEncodedBodyHash() string {
	if val, ok := s.Request.Header[s.Constant.ContentSha256Header]; ok {
		return val[0]
	}
	if s.Request.Body == nil {
		return s.hexEncodedHash("")
	}
	bodyBytes, _ := ioutil.ReadAll(s.Request.Body)
	return s.hexEncodedHash(string(bodyBytes))
}

func (s *Signer) isSignableHeader(key string) bool {
	nonSignableHeaders := []string{"authorization", "content-type", "content-length", "user-agent", "presigned-expires", "expect", "x-amzn-trace-id"}
	if strings.HasPrefix(key, "x-amz-") {
		return true
	}
	for _, h := range nonSignableHeaders {
		if h == key {
			return false
		}
	}
	return true
}

func (s *Signer) iso8601(date time.Time) string {
	return date.UTC().Format("20060102T150405Z")
}

func (s *Signer) getSigningKey(credentials map[string]string, date, region, service string) []byte {
	kDate := hmacSHA256([]byte(s.Constant.KDatePrefix+credentials["secretAccessKey"]), []byte(date))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte(s.Constant.V4Identifier))
	return kSigning
}

func (s *Signer) createScope(date, region, service string) string {
	return strings.Join([]string{date, region, service, s.Constant.V4Identifier}, "/")
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func main() {
	// Example usage
	req, _ := http.NewRequest("GET", "http://example.amazonaws.com/?Param2=value2&Param1=value1", nil)
	req.Header.Set("X-Amz-Date", "20150830T123600Z")
	req.Header.Set("Host", "example.amazonaws.com")

	credentials := map[string]string{
		"accessKeyId":     "AKIDEXAMPLE",
		"secretAccessKey": "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
	}

	signer := NewSigner(req, "service", "us-east-1", false)
	date, _ := time.Parse("20060102T150405Z", "20150830T123600Z")
	signer.AddAuthorization(credentials, date)

	fmt.Println(req.Header.Get("Authorization"))
}
