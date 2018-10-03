package square_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/timhugh/square"
)

func ExampleAuthenticateRequest() {
	url := "http://www.example.com/events"
	body := `{"event": "test"}`
	signature := "n96t75ZEk8OvwpqHZk/O4HMnt1E="

	// Creating a request like one we would receive from a *http.Handler
	request, _ := http.NewRequest("POST", url, strings.NewReader(body))
	request.Header.Set("X-Square-Signature", signature)

	signatureKey := "example_key"

	if err := square.AuthenticateRequest(request, signatureKey); err != nil {
		fmt.Print("Bad request!")
	} else {
		fmt.Print("Success")
	}

	// Output: Success
}

func ExampleAuthenticateSignature() {
	url := "http://www.example.com/events"
	body := `{"event": "test"}`
	signature := "n96t75ZEk8OvwpqHZk/O4HMnt1E="

	signatureKey := "example_key"

	if err := square.AuthenticateSignature(signature, url, body, signatureKey); err != nil {
		fmt.Print("Bad request!")
	} else {
		fmt.Print("Success")
	}

	// Output: Success
}

func ExampleGenerateSignature() {
	url := "http://www.example.com/events"
	body := `{"event": "test"}`
	key := "example_key"

	signature := square.GenerateSignature(url, body, key)
	fmt.Print(signature)

	// Output: n96t75ZEk8OvwpqHZk/O4HMnt1E=
}
