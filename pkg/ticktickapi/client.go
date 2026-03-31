package ticktickapi

import (
	"fmt"

	"github.com/alex-305/ticktui/internal/config"
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

func GetClient() (*Client, error) {
	token, err := config.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	return NewClient(token), nil
}
