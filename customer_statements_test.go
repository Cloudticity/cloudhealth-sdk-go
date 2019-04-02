package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultBillingArtifact = BillingArtifact{
	CustomerID:                           1234567890,
	BillingPeriod:                        "1970-01-01",
	TotalAmount:                          9999.99,
	Status:                               "Final",
	DetailedBillingRecordsGenerationTime: "1970-01-01",
	StatementGenerationTime:              "1970-01-01",
	Currency:                             Currency{Name: "USD", Symbol: "$"},
}

var defaultBillingArtifacts = BillingArtifacts{
	BillingArtifacts: []BillingArtifact{
		{
			CustomerID:                           1234567890,
			BillingPeriod:                        "1970-01-01",
			TotalAmount:                          9999.99,
			Status:                               "Final",
			DetailedBillingRecordsGenerationTime: "1970-01-01",
			StatementGenerationTime:              "1970-01-01",
			Currency:                             Currency{Name: "USD", Symbol: "$"},
		},
		{
			CustomerID:                           98765433210,
			BillingPeriod:                        "1970-01-01",
			TotalAmount:                          110.00,
			Status:                               "Final",
			DetailedBillingRecordsGenerationTime: "1970-01-01",
			StatementGenerationTime:              "1970-01-01",
			Currency:                             Currency{Name: "USD", Symbol: "$"},
		},
	},
}

func TestGetSingleBillingArtifact(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/customer_statements/%d", defaultBillingArtifact.CustomerID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultBillingArtifact)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedBillingArtifact, err := c.GetSingleBillingArtifacts(1234567890)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if returnedBillingArtifact.CustomerID != defaultBillingArtifact.CustomerID {
		t.Errorf("GetBillingArtifact() expected CustomerID `%d`, got `%d`", defaultBillingArtifact.CustomerID, returnedBillingArtifact.CustomerID)
		return
	}
}

func TestGetBillingArtifacts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := "/customer_statements/"
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultBillingArtifacts)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedBillingArtifacts, err := c.GetBillingArtifacts()
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if len(returnedBillingArtifacts.BillingArtifacts) != 2 {
		t.Errorf("All customer statements have not been retrieved")
		return
	}
	if returnedBillingArtifacts.BillingArtifacts[0].CustomerID != defaultBillingArtifacts.BillingArtifacts[0].CustomerID && returnedBillingArtifacts.BillingArtifacts[1].CustomerID != defaultBillingArtifacts.BillingArtifacts[1].CustomerID {
		t.Errorf("Retrieved customer statements don't match")
		return
	}
	return
}
