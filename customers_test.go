package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultCustomer = Customer{
	ID:   1234567890,
	Name: "test",
}

var defaultCustomers = Customers{
	Customers: []Customer{
		Customer{
			ID:   1234567890,
			Name: "test",
		},
		Customer{
			ID:   9876543210,
			Name: "tset",
		},
	},
}

func TestGetCustomers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := "/customers/"
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultCustomers)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCustomers, err := c.GetCustomers()
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if len(returnedCustomers.Customers) != 2 {
		t.Errorf("All customers have not been retrieved")
		return
	}
	if returnedCustomers.Customers[0].ID != defaultCustomers.Customers[0].ID && returnedCustomers.Customers[0].ID != defaultCustomers.Customers[1].ID {
		t.Errorf("Retrieved customers don't match")
		return
	}
	return
}

func TestGetSingleCustomer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customers/%d", defaultAWSAccount.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultCustomer)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCustomer, err := c.GetSingleCustomer(1234567890)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if returnedCustomer.ID != defaultCustomer.ID {
		t.Errorf("GetCustomer() expected ID `%d`, got `%d`", defaultCustomer.ID, returnedCustomer.ID)
		return
	}
}

func TestGetSingleCustomerDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customers/%d", defaultAWSAccount.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultCustomer)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	_, err = c.GetSingleCustomer(defaultCustomer.ID)
	if err != ErrCustomerNotFound {
		t.Errorf("GetAwsAccount() returned the wrong error: %s", err)
		return
	}
}
