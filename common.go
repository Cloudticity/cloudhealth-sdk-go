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

// ErrHeaderMissing is returned when an header is missing (400).
var ErrHeaderMissing = errors.New("Header missing")

// ErrNotFound is returned when a Resource doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if not exists) (404).
var ErrNotFound = errors.New("Resource not found")

// ErrClientAuthenticationError is returned for authentication errors with the API (401).
var ErrClientAuthenticationError = errors.New("Authentication Error with CloudHealth")

// ErrForbidden is returned for access rights errors with the API (403).
var ErrForbidden = errors.New("Authentication Error with CloudHealth")

// ErrUnprocessableEntityError is returned for resource creation errors (422).
var ErrUnprocessableEntityError = errors.New("Bad Request (Input format error). Please check if a resource with same name already exists")

// ErrTooManyRequest is returned by API when it's throttling you (429).
var ErrTooManyRequest = errors.New("Exceeding post rate limit")

// getResponsePage returns a response page from a CloudHealth's endpoint.
func getResponsePage(s *Client, relativeURL *url.URL) ([]byte, error) {
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return []byte{}, err
	}

	return sendRequest(req)
}

// createResource creates a resource and retrieves details from CloudHealth.
func createResource(s *Client, relativeURL *url.URL, resource interface{}) ([]byte, error) {
	body, _ := json.Marshal(resource)
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	return sendRequest(req)
}

// updateResource updates a resource and retrieves details from CloudHealth.
func updateResource(s *Client, relativeURL *url.URL, resource interface{}) ([]byte, error) {
	body, _ := json.Marshal(resource)
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	return sendRequest(req)
}

// deleteResource deletes a resource and retrieves details from CloudHealth.
func deleteResource(s *Client, relativeURL *url.URL) ([]byte, error) {
	url := s.EndpointURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	return sendRequest(req)
}

// sendRequest sends request to CloudHealth and retrieves details about.
func sendRequest(req *http.Request) ([]byte, error) {
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
	case http.StatusOK, http.StatusCreated, http.StatusNoContent: //200, 201, 204
		return responseBody, nil
	case http.StatusBadRequest: //400
		return nil, ErrHeaderMissing
	case http.StatusUnauthorized: //401
		return nil, ErrClientAuthenticationError
	case http.StatusForbidden: //403
		return nil, ErrForbidden
	case http.StatusNotFound: //404
		return nil, ErrNotFound
	case http.StatusUnprocessableEntity: //422
		return nil, ErrUnprocessableEntityError
	case http.StatusTooManyRequests: //429
		return nil, ErrTooManyRequest
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}
