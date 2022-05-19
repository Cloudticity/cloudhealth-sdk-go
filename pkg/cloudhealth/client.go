// Package cloudhealth is a wrapper for the CloudHealth API.
package cloudhealth

// Client communicates with the CloudHealth API.
type Client struct {
	APIKey      string
	EndpointURL string
}

// NewClient returns a new CloudHealth.Client for accessing the CloudHealth API.
func NewClient(apiKey string, defaultEndpointURL string) (*Client, error) {
	s := &Client{
		APIKey: apiKey,
	}

	s.EndpointURL = defaultEndpointURL
	return s, nil
}
