package criblclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var DefaultRestURL = "http://" + os.Getenv("CRIBL_HOST") + "/api/v1"

type Client struct {
	Host       string
	HTTPClient *http.Client
	Token      string
	Auth       Auth
}

// func NewClient(apiKey string) *Client {
// 	return &Client{
// 		HttpClient: http.DefaultClient,
// 		ApiKey:     apiKey,
// 	}
// }

func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Host:       DefaultRestURL,
	}

	if host != nil {
		c.Host = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = Auth{
		Username: *username,
		Password: *password,
	}

	ar, err := c.AuthLogin()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) newRequest(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", DefaultRestURL, path), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		return body, err
	} else {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}
}
