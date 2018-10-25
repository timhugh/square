package square_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/timhugh/square"
)

func TestFetchesPayments(t *testing.T) {
	server, client := buildClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("expected request to include authorization header \"Bearer token\" but got \"%s\"", r.Header.Get("Authorization"))
		}
		fmt.Fprint(w, `{"payment_id": "payment_id", "location_id": "location_id"}`)
	}))
	defer server.Close()

	paymentData, err := client.FetchPayment("token", "payment_id", "location_id")
	if err != nil {
		t.Errorf("expected to retrieve payment without error but got \"%s\"", err)
	}

	expectedData := map[string]interface{}{
		"payment_id":  "payment_id",
		"location_id": "location_id",
	}
	if !reflect.DeepEqual(expectedData, paymentData) {
		t.Errorf("expected to receive:\n%+v\ngot:\n%+v\n", expectedData, paymentData)
	}
}

func TestFetchPaymentNotFound(t *testing.T) {
	server, client := buildClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"type":"not_found","message":"Payment not found"}`)
	}))
	defer server.Close()

	_, err := client.FetchPayment("token", "payment_id", "location_id")

	_, ok := err.(NotFoundError)
	if err == nil || err.Error() != "Payment not found" || !ok {
		t.Errorf("expected to receive NotFoundError with message \"Payment not found\" error but got %T with message \"%s\"", err, err)
	}
}

func TestFetchPaymentAuthError(t *testing.T) {
	server, client := buildClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"type":"service.not_authorized","message":"Not Authorized"}`)
	}))
	defer server.Close()

	_, err := client.FetchPayment("token", "payment_id", "location_id")

	_, ok := err.(NotAuthorizedError)
	if err == nil || err.Error() != "Not Authorized" || !ok {
		t.Errorf("expected to receive NotAuthorizedError with message \"Not Authorized\" error but got %T with message \"%s\"", err, err)
	}
}

func buildClient(handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewTLSServer(handler)
	client := &Client{
		ApiUrl:     server.URL,
		HttpClient: server.Client(),
	}
	return server, client
}
