package cloudhealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ErrNotFound is returned when a Resource doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if exists).
var ErrNotFound = errors.New("Resource not found")

// getResponsePage returns a response page of a CloudHealth's endpoint.
func getResponsePage(s *Client, relativeURL *url.URL) ([]byte, error) {
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("GET", url.String(), nil)

	client := &http.Client{
		Timeout: time.Second * 20,
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
		return nil, ErrNotFound
	default:
		return []byte{}, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// createResource creates a resource and retrieves details from CloudHealth.
func createResource(s *Client, relativeURL *url.URL, resource interface{}) ([]byte, error) {
	body, _ := json.Marshal(resource)
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 20,
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
	case http.StatusCreated:
		return responseBody, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntityError
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// updateResource updates a resource and retrieves details from CloudHealth.
func updateResource(s *Client, relativeURL *url.URL, resource interface{}) ([]byte, error) {
	body, _ := json.Marshal(resource)
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 20,
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
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntityError
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}
