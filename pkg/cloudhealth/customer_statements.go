package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// BillingArtifacts represents all Customer Statements enabled in CloudHealth with their details.
type BillingArtifacts struct {
	BillingArtifacts []BillingArtifact `json:"billing_artifacts"`
}

// BillingArtifact represents the configuration of a Customer Statement in CloudHealth with its details.
type BillingArtifact struct {
	CustomerID                           int      `json:"customer_id"`
	CloudProvider                        string   `json:"cloud"`
	BillingPeriod                        string   `json:"billing_period"`
	TotalAmount                          float64  `json:"total_amount"`
	Status                               string   `json:"status"`
	DetailedBillingRecordsGenerationTime string   `json:"detailed_billing_records_generation_time"`
	StatementGenerationTime              string   `json:"statement_generation_time"`
	StatementSummaryGenerationTime       string   `json:"statement_summary_generation_time"`
	Currency                             Currency `json:"currency"`
	InvoiceId                            string   `json:"invoice_id,omitempty"`
	InvoiceDate                          string   `json:"invoice_date,omitempty"`
}

// Currency represents the currency used for billing.
type Currency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// GetSingleCustomerStatements gets all statements for a specific Customer ID.
func (s *Client) GetSingleCustomerStatements(id int) (*BillingArtifacts, error) {
	// Set variables we will need along the way
	var billingArtifacts BillingArtifacts
	var page, pageSize int = 1, 100

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v1/customer_statements?client_api_id=%d&%s", id, params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the BillingArtifacts struct
		err = json.Unmarshal(responseBody, &billingArtifacts)
		if err != nil {
			return nil, err
		}

		// Check length of array in BillingArtifacts' struct to determine if we should break out of the loop
		if len(billingArtifacts.BillingArtifacts) < pageSize {
			break
		}

		// Increment page counter
		page++
	}

	return &billingArtifacts, nil
}

// GetCustomerStatements gets all Statements.
func (s *Client) GetCustomerStatements() (*BillingArtifacts, error) {
	// Set variables we will need along the way
	var billingArtifacts BillingArtifacts
	var page, pageSize int = 1, 100

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v1/customer_statements?%s", params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the BillingArtifacts struct
		err = json.Unmarshal(responseBody, &billingArtifacts)
		if err != nil {
			return nil, err
		}

		// Check length of array in BillingArtifacts' struct to determine if we should break out of the loop
		if len(billingArtifacts.BillingArtifacts) < pageSize {
			break
		}

		// Increment page counter
		page++
	}

	return &billingArtifacts, nil
}
