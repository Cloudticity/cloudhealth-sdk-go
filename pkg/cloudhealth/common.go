package cloudhealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ErrHeaderMissing is returned when an header is missing (400).
var ErrHeaderMissing = errors.New("header missing")

// ErrNotFound is returned when a Resource doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if not exists) (404).
var ErrNotFound = errors.New("resource not found")

// ErrClientAuthenticationError is returned for authentication errors with the API (401).
var ErrClientAuthenticationError = errors.New("authentication error with CloudHealth")

// ErrForbidden is returned for access rights errors with the API (403).
var ErrForbidden = errors.New("authentication error with CloudHealth")

// ErrUnprocessableEntityError is returned for resource creation errors (422).
var ErrUnprocessableEntityError = errors.New("bad Request (Input format error). Please check if a resource with same name already exists")

// ErrTooManyRequest is returned by API when it's throttling you (429).
var ErrTooManyRequest = errors.New("exceeding post rate limit")

// getResponsePage returns a response page from a CloudHealth's endpoint.
func getResponsePage(s *Client, relativeURL string) ([]byte, error) {
	// Set up the URL
	finalUrl := s.EndpointURL + relativeURL

	// Make the physical API call
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	// Add headers as needed
	req.Header.Add("Authorization", s.APIKey)
	req.Header.Add("Accept", "application/json")

	return sendRequest(req)
}

// createResource creates a resource and retrieves details from CloudHealth.
func createResource(s *Client, relativeURL string, resource interface{}) ([]byte, error) {
	// Create the request body
	body, _ := json.Marshal(resource)

	// Set up the URL
	finalUrl := s.EndpointURL + relativeURL

	// Make the physical API call
	req, err := http.NewRequest("POST", finalUrl, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}

	// Add headers as needed
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.APIKey)

	return sendRequest(req)
}

// updateResource updates a resource and retrieves details from CloudHealth.
func updateResource(s *Client, relativeURL string, resource interface{}) ([]byte, error) {
	// Create the request body
	body, _ := json.Marshal(resource)

	// Set up the URL
	finalUrl := s.EndpointURL + relativeURL

	// Make the physical API call
	req, err := http.NewRequest("PUT", finalUrl, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}

	// Add headers as needed
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.APIKey)

	return sendRequest(req)
}

// deleteResource deletes a resource and retrieves details from CloudHealth.
func deleteResource(s *Client, relativeURL string) ([]byte, error) {
	// Set up the URL
	finalUrl := s.EndpointURL + relativeURL

	// Make the physical API call
	req, err := http.NewRequest("DELETE", finalUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	// Add headers as needed
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.APIKey)

	return sendRequest(req)
}

// sendRequest sends request to CloudHealth and retrieves details about.
func sendRequest(req *http.Request) ([]byte, error) {
	// Create HTTP client for sending requests
	client := &http.Client{
		Timeout: time.Second * 20,
	}

	// Make the API call
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read in the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check and handle the HTTP status code of the return
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
		return nil, fmt.Errorf("unknown status code response from CloudHealth: `%d`", resp.StatusCode)
	}
}
