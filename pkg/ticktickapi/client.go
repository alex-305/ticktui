package ticktickapi

import (
	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "https://api.ticktick.com/open/v1"
)

type Client struct {
	http *resty.Client
}

func NewClient(token string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Authorization", "Bearer "+token)

	return &Client{http: client}
}

func GetClient(token string) (*Client, error) {
	return NewClient(token), nil
}
