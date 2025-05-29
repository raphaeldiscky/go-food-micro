// Package client provides a http client.
package client

import (
	"time"

	resty "github.com/go-resty/resty/v2"
)

// Constants for the http client.
const (
	timeout               = 5 * time.Second
	dialContextTimeout    = 5 * time.Second
	tLSHandshakeTimeout   = 5 * time.Second
	xaxIdleConns          = 20
	maxConnsPerHost       = 40
	retryCount            = 3
	retryWaitTime         = 300 * time.Millisecond
	idleConnTimeout       = 120 * time.Second
	responseHeaderTimeout = 5 * time.Second
)

// NewHttpClient is a function that creates a new http client.
func NewHttpClient() *resty.Client {
	client := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime)

	return client
}
