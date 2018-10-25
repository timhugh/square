package square

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	ApiUrl     string
	HttpClient *http.Client
}

const squareApiUrl = "https://connect.squareup.com"
const (
	paymentRoute = "/v1/%s/payments/%s"
)

func NewClient() *Client {
	return &Client{
		ApiUrl:     squareApiUrl,
		HttpClient: http.DefaultClient,
	}
}

func (c *Client) FetchPayment(authToken, paymentID, locationID string) (map[string]interface{}, error) {
	url := c.ApiUrl + fmt.Sprintf(paymentRoute, locationID, paymentID)
	return c.getJSONResponse(url, authToken)
}

func (c *Client) getJSONResponse(url, authToken string) (map[string]interface{}, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := readBody(response.Body)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, NotFoundError(jsonData["message"].(string))
	}
	if response.StatusCode == http.StatusUnauthorized {
		return nil, NotAuthorizedError(jsonData["message"].(string))
	}
	return jsonData, nil
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

type NotAuthorizedError string

func (e NotAuthorizedError) Error() string {
	return string(e)
}
