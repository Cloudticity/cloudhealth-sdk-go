package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultCustomerStatement = CustomerStatement{
	CustomerID:    1234567890,
	BillingPeriod: "1970-01-01",
	TotalAmount:   9999.99,
	Status:        "Final",
	DetailedBillingRecordsGenerationTime: "1970-01-01",
	StatementGenerationTime:              "1970-01-01",
	Currency:                             Currency{Name: "USD", Symbol: "$"},
}

var defaultCustomerStatements = CustomerStatements{
	CustomerStatements: []CustomerStatement{
		{
			CustomerID:    1234567890,
			BillingPeriod: "1970-01-01",
			TotalAmount:   9999.99,
			Status:        "Final",
			DetailedBillingRecordsGenerationTime: "1970-01-01",
			StatementGenerationTime:              "1970-01-01",
			Currency:                             Currency{Name: "USD", Symbol: "$"},
		},
		{
			CustomerID:    98765433210,
			BillingPeriod: "1970-01-01",
			TotalAmount:   110.00,
			Status:        "Final",
			DetailedBillingRecordsGenerationTime: "1970-01-01",
			StatementGenerationTime:              "1970-01-01",
			Currency:                             Currency{Name: "USD", Symbol: "$"},
		},
	},
}

func TestGetSingleCustomerStatement(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customer_statements/%d", defaultCustomerStatement.CustomerID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultCustomerStatement)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCustomerStatement, err := c.GetSingleCustomerStatement(1234567890)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if returnedCustomerStatement.CustomerID != defaultCustomerStatement.CustomerID {
		t.Errorf("GetCustomerStatement() expected CustomerID `%d`, got `%d`", defaultCustomerStatement.CustomerID, returnedCustomerStatement.CustomerID)
		return
	}
}

func TestGetCustomerStatements(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := "/customer_statements/"
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultCustomerStatements)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCustomerStatements, err := c.GetCustomerStatements()
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if len(returnedCustomerStatements.CustomerStatements) != 2 {
		t.Errorf("All customer statements have not been retrieved")
		return
	}
	if returnedCustomerStatements.CustomerStatements[0].CustomerID != defaultCustomerStatements.CustomerStatements[0].CustomerID && returnedCustomerStatements.CustomerStatements[1].CustomerID != defaultCustomerStatements.CustomerStatements[1].CustomerID {
		t.Errorf("Retrieved customer statements don't match")
		return
	}
	return
}
