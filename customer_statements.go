package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// CustomerStatements represents all Customer Statements enabled in CloudHealth with their details.
type CustomerStatements struct {
	CustomerStatements []CustomerStatement `json:"billing_artifacts"`
}

// CustomerStatement represents the configuration of a Customer Statement in CloudHealth with its details.
type CustomerStatement struct {
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
func (s *Client) GetSingleCustomerStatement(id int) (*CustomerStatement, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/%d?api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var billingArtifact = new(CustomerStatement)
	err = json.Unmarshal(responseBody, &billingArtifact)
	if err != nil {
		return nil, err
	}

	return billingArtifact, nil
}

// GetCustomerStatements gets all Customer Statements.
func (s *Client) GetCustomerStatements() (*CustomerStatements, error) {
	customerStatements := new(CustomerStatements)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"100"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var ba = new(CustomerStatements)
		err = json.Unmarshal(responseBody, &ba)
		if err != nil {
			return nil, err
		}
		for _, a := range ba.CustomerStatements {
			customerStatements.CustomerStatements = append(customerStatements.CustomerStatements, a)
		}
		if len(ba.CustomerStatements) < 100 {
			break
		}
		page++
	}

	return customerStatements, nil
}
