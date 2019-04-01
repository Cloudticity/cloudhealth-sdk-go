// Package cloudhealth is a wrapper for the CloudHealth API.
package cloudhealth

import (
	"net/url"
)

// Client communicates with the CloudHealth API.
type Client struct {
	APIKey      string
	EndpointURL *url.URL
}

// NewClient returns a new cloudhealth.Client for accessing the CloudHealth API.
func NewClient(apiKey string, defaultEndpointURL string) (*Client, error) {
	s := &Client{
		APIKey: apiKey,
	}
	endpointURL, err := url.Parse(defaultEndpointURL)
	if err != nil {
		return nil, err
	}
	s.EndpointURL = endpointURL
	return s, nil
}
