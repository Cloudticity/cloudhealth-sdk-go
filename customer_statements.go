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
	var billingArtifacts = new(BillingArtifacts)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"100"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/?client_api_id=%d&%s", id, params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var ba = new(BillingArtifacts)
		err = json.Unmarshal(responseBody, &ba)
		if err != nil {
			return nil, err
		}
		for _, a := range ba.BillingArtifacts {
			billingArtifacts.BillingArtifacts = append(billingArtifacts.BillingArtifacts, a)
		}
		if len(ba.BillingArtifacts) < 100 {
			break
		}
		page++
	}

	return billingArtifacts, nil
}

// GetCustomerStatements gets all Statements.
func (s *Client) GetCustomerStatements() (*BillingArtifacts, error) {
	billingArtifacts := new(BillingArtifacts)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"100"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var ba = new(BillingArtifacts)
		err = json.Unmarshal(responseBody, &ba)
		if err != nil {
			return nil, err
		}
		for _, a := range ba.BillingArtifacts {
			billingArtifacts.BillingArtifacts = append(billingArtifacts.BillingArtifacts, a)
		}
		if len(ba.BillingArtifacts) < 100 {
			break
		}
		page++
	}

	return billingArtifacts, nil
}
