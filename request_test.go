package square_test

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/timhugh/square"
)

const goodSignature = "DwwpgL5sy1WXHwPSsLNN27tGRSY="
const requestBody = `{"event": "test"}`
const requestURL = "http://www.example.com/events"

const signatureKey = "test_key"

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
	if err == nil {
		t.Fatal("expected InvalidSignature error but got none")
	}

	_, ok := err.(square.InvalidSignature)
	if !ok {
		t.Errorf("expected InvalidSignature error but got %+v", err)
	}
}

func TestGenerateSignature(t *testing.T) {
	signature := square.GenerateSignature(requestURL, requestBody, signatureKey)
	if signature != goodSignature {
		t.Errorf("expected %s, got %s", goodSignature, signature)
	}
}
