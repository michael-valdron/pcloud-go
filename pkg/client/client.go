package client

import "net/http"

type Client struct {
	Auth   *string
	Client *http.Client
}

// NewClient create new pCloudClient
func NewClient() *Client {
	return &Client{
		Auth:   nil,
		Client: &http.Client{},
	}
}
