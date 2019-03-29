package cloudhealth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// getResponsePage returns a response page of a CloudHealth's endpoint.
func getResponsePage(s *Client, relativeURL *url.URL) ([]byte, error) {
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("GET", url.String(), nil)

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return responseBody, nil
	case http.StatusForbidden:
		return []byte{}, ErrClientAuthenticationError
	case http.StatusNotFound:
		return nil, ErrAwsAccountNotFound
	default:
		return []byte{}, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}
