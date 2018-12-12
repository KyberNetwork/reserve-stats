package core

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client is the the real implementation of core client interface.
type Client struct {
	sugar      *zap.SugaredLogger
	client     *http.Client
	url        string
	signingKey string
}

type commonResponse struct {
	Reason  string `json:"reason"`
	Success bool   `json:"success"`
}

// NewClient creates a new core client instance.
func NewClient(sugar *zap.SugaredLogger, url, signingKey string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar, url: url, signingKey: signingKey, client: client}, nil
}
