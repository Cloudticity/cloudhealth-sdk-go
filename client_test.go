package cloudhealth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadAPIKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/aws_accounts/")
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	_, err = c.GetAwsAccounts()
	if err != ErrForbidden {
		t.Errorf("GetAwsAccounts() returned the wrong error: %s", err)
		return
	}
}

func TestBadEndPoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	_, err = c.GetAwsAccounts()
	if err != ErrClientAuthenticationError {
		t.Errorf("GetAwsAccounts() returned the wrong error: %s", err)
		return
	}
}
