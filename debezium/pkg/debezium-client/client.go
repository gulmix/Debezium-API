package debezium_client

import (
	"net/http"
	"strings"
	"time"
)

type Client struct {
	cc      *http.Client
	baseURL string
}

func New(baseUrl string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseUrl, "/"),
		cc:      &http.Client{Timeout: timeout},
	}
}
