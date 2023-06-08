package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client  *resty.Client
	baseURL string
}

func New(url string) *Client {
	hc := &http.Client{}
	rc := resty.NewWithClient(hc)

	return &Client{
		client:  rc,
		baseURL: url,
	}
}

func (c *Client) path(f string, args ...any) string {
	return fmt.Sprintf(c.baseURL+f, args...)
}
