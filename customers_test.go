package cloudhealth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultCustomer = Customer{
	ID:   1234567890,
	Name: "test",
	Tags: map[string]string{"key": "value"},
}

var defaultCustomers = Customers{
	Customers: []Customer{
		{
			ID:   1234567890,
			Name: "test",
			Tags: map[string]string{"key": "value"},
		},
		{
			ID:   9876543210,
			Name: "tset",
			Tags: map[string]string{"key": "value"},
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
	if returnedCustomers.Customers[0].ID != defaultCustomers.Customers[0].ID && returnedCustomers.Customers[1].ID != defaultCustomers.Customers[1].ID {
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
		expectedURL := fmt.Sprintf("/customers/%d", defaultCustomer.ID)
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
		expectedURL := fmt.Sprintf("/customers/121212121")
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

	_, err = c.GetSingleCustomer(121212121)
	if err != ErrNotFound {
		t.Errorf("GetCustomer() returned the wrong error: %s", err)
		return
	}
}

func TestCreateCustomer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if r.URL.EscapedPath() != "/customers" {
			t.Errorf("Expected request to ‘/customers, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("Unable to read response body")
		}

		customer := new(Customer)
		err = json.Unmarshal(body, &customer)
		if err != nil {
			t.Errorf("Unable to unmarshal Customer, got `%s`", body)
		}
		if customer.Name != "test" {
			t.Errorf("Expected request to include customer name ‘test’, got ‘%s’", customer.Name)
		}
		customer.ID = 1234567890
		js, err := json.Marshal(customer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCustomer, err := c.CreateCustomer(Customer{
		Name: "test",
	})
	if err != nil {
		t.Errorf("CreateCustomer() returned an error: %s", err)
		return
	}
	if returnedCustomer.Name != "test" {
		t.Errorf("Createcustomer() expected name `test`, got `%s`", returnedCustomer.Name)
		return
	}
}

func TestUpdateCustomer(t *testing.T) {
	updatedCustomer := defaultCustomer
	updatedCustomer.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customers/%d", defaultCustomer.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(updatedCustomer)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCustomer, err := c.UpdateCustomer(updatedCustomer)
	if err != nil {
		t.Errorf("UpdateCustomer() returned an error: %s", err)
		return
	}
	if returnedCustomer.ID != updatedCustomer.ID {
		t.Errorf("UpdateCustomer() expected ID `%d`, got `%d`", defaultCustomer.ID, returnedCustomer.ID)
		return
	}
	if returnedCustomer.Name == defaultCustomer.Name {
		t.Errorf("UpdateCustomer() did not update the name")
		return
	}
}

func TestDeleteCustomer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customers/%d", defaultCustomer.ID)
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

	err = c.DeleteCustomer(defaultCustomer.ID)
	if err != nil {
		t.Errorf("DeleteCustomer() returned an error: %s", err)
		return
	}
}

func TestDeleteCustomerDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customers/%d", defaultCustomer.ID)
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

	err = c.DeleteCustomer(defaultCustomer.ID)
	if err != ErrNotFound {
		t.Errorf("DeleteCustomer() returned the wrong error: %s", err)
		return
	}
}
