package square_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/timhugh/square"
)

const (
	goodSignature = "DwwpgL5sy1WXHwPSsLNN27tGRSY="
	requestBody   = `{"event": "test"}`
	requestURL    = "http://www.example.com/events"
	signatureKey  = "test_key"
)

func stubRequest(url, body, signature string) *http.Request {
	r := httptest.NewRequest("POST", requestURL, strings.NewReader(requestBody))
	r.Header.Set("X-Square-Signature", signature)
	return r
}

func TestGoodSignature(t *testing.T) {
	request := stubRequest(requestURL, requestBody, goodSignature)

	err := square.AuthenticateRequest(request, signatureKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBadSignature(t *testing.T) {
	request := stubRequest(requestURL, requestBody, "bad_signature")

	err := square.AuthenticateRequest(request, signatureKey)
	_, ok := err.(square.InvalidSignatureError)
	if err == nil || !strings.Contains(err.Error(), "invalid signature") || !ok {
		t.Errorf("expected InvalidSignatureError error but got %+v", err)
	}
}

func TestBodyReadError(t *testing.T) {
	body := &ErrReader{}
	request := httptest.NewRequest("POST", requestURL, body)
	request.Header.Set("X-Square-Signature", goodSignature)

	err := square.AuthenticateRequest(request, signatureKey)
	if err == nil || !strings.Contains(err.Error(), "read error") {
		t.Fatalf("expected read error but got %+v", err)
	}
}

func TestGenerateSignature(t *testing.T) {
	signature := square.GenerateSignature(requestURL, requestBody, signatureKey)
	if signature != goodSignature {
		t.Errorf("expected %s, got %s", goodSignature, signature)
	}
}

type ErrReader struct{}

func (r *ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
