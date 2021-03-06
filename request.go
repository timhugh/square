package square

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// InvalidSignatureError is used to report an error when a signature cannot be authenticated.
type InvalidSignatureError string

func (r InvalidSignatureError) Error() string {
	return "invalid signature" + string(r)
}

// AuthenticateRequest authenticates an HTTP request against a signature key by comparing the X-Square-Signature header to a valid signature generated from the request url, body, and the signature key.
// A valid request will return nil. An invalid signature will return a InvalidSignatureError error.
// Any other errors will be directly returned.
func AuthenticateRequest(r *http.Request, signatureKey string) error {
	requestURL := r.URL.String()
	requestBody, err := readBody(r.Body)
	// replace the request body so it can be read again later
	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	requestSignature := r.Header.Get("X-Square-Signature")

	return AuthenticateSignature(requestSignature, requestURL, string(requestBody), signatureKey)
}

// AuthenticateSignature authenticates a signature against a signature key by comparing the signature to a valid signature generated from the url, body, and the signature key.
// A valid signature will return nil. An invalid signature will return an InvalidSignatureError error.
func AuthenticateSignature(signature, url, body, signatureKey string) error {
	expectedSignature := GenerateSignature(url, body, signatureKey)
	if signature != expectedSignature {
		return InvalidSignatureError(fmt.Sprintf("expected \"%s\", got \"%s\"", expectedSignature, signature))
	}
	return nil
}

// GenerateSignature creates a valid Base64-encoded HMAC-SHA1 signature from a request url, body, and a signature key.
func GenerateSignature(url, body, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(url + body))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func readBody(body io.ReadCloser) ([]byte, error) {
	buf, err := ioutil.ReadAll(io.LimitReader(body, 1048576))
	body.Close()
	if err != nil {
		return nil, err
	}
	return buf, nil
}
