package sdk

import "net/url"

// SDK defines the interface for all platform SDKs.
type SDK interface {
	// Sign generates a signature for a request.
	Sign(method, uri, a1, xsecAppid string, params url.Values) string
}
