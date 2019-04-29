package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// BillingStatements represents all Customer Statements enabled in CloudHealth with their details.
type BillingStatements struct {
	BillingStatements []BillingStatement `json:"billing_artifacts"`
}

// BillingStatement represents the configuration of a Customer Statement in CloudHealth with its details.
type BillingStatement struct {
	CustomerID                           int      `json:"customer_id"`
	BillingPeriod                        string   `json:"billing_period"`
	TotalAmount                          float64  `json:"total_amount"`
	Status                               string   `json:"status"`
	DetailedBillingRecordsGenerationTime string   `json:"detailed_billing_records_generation_time"`
	StatementGenerationTime              string   `json:"statement_generation_time"`
	StatementSummaryGenerationTime       string   `json:"statement_summary_generation_time"`
	Currency                             Currency `json:"currency"`
}

// Currency represents the currency used for billing.
type Currency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// GetSingleCustomerStatement gets details for the specified CloudHealth Customer Statement ID.
func (s *Client) GetSingleCustomerStatement(id int) (*BillingStatement, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/?client_api_id=%d&api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var billingArtifact = new(BillingStatement)
	err = json.Unmarshal(responseBody, &billingArtifact)
	if err != nil {
		return nil, err
	}

	return billingArtifact, nil
}

// GetCustomerStatements gets all Statements.
func (s *Client) GetCustomerStatements() (*BillingStatements, error) {
	billingStatements := new(BillingStatements)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"100"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var ba = new(BillingStatements)
		err = json.Unmarshal(responseBody, &ba)
		if err != nil {
			return nil, err
		}
		for _, a := range ba.BillingStatements {
			billingStatements.BillingStatements = append(billingStatements.BillingStatements, a)
		}
		if len(ba.BillingStatements) < 100 {
			break
		}
		page++
	}

	return billingStatements, nil
}
