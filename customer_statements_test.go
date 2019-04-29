package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultBillingStatement = BillingStatement{
	CustomerID:    1234567890,
	BillingPeriod: "1970-01-01",
	TotalAmount:   9999.99,
	Status:        "Final",
	DetailedBillingRecordsGenerationTime: "1970-01-01",
	StatementGenerationTime:              "1970-01-01",
	Currency:                             Currency{Name: "USD", Symbol: "$"},
}

var defaultBillingStatements = BillingStatements{
	BillingStatements: []BillingStatement{
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

func TestGetSingleBillingStatement(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customer_statements/")
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultBillingStatement)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedBillingStatement, err := c.GetSingleCustomerStatement(1234567890)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if returnedBillingStatement.CustomerID != defaultBillingStatement.CustomerID {
		t.Errorf("GetBillingStatement() expected CustomerID `%d`, got `%d`", defaultBillingStatement.CustomerID, returnedBillingStatement.CustomerID)
		return
	}
}

func TestGetBillingStatements(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := "/customer_statements/"
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultBillingStatements)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedBillingStatements, err := c.GetCustomerStatements()
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if len(returnedBillingStatements.BillingStatements) != 2 {
		t.Errorf("All customer statements have not been retrieved")
		return
	}
	if returnedBillingStatements.BillingStatements[0].CustomerID != defaultBillingStatements.BillingStatements[0].CustomerID && returnedBillingStatements.BillingStatements[1].CustomerID != defaultBillingStatements.BillingStatements[1].CustomerID {
		t.Errorf("Retrieved customer statements don't match")
		return
	}
	return
}
