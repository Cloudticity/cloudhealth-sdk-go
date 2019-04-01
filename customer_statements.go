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
	CustomerID                           int     `json:"customer_id"`
	BillingPeriod                        string  `json:"billing_period"`
	TotalAmount                          float64 `json:"total_amount"`
	Status                               string  `json:"status"`
	DetailedBillingRecordsGenerationTime string  `json:"detailed_billing_records_generation_time"`
	StatementGenerationTime              string  `json:"statement_generation_time"`
}

// Currency represents the currency used for billing.
type Currency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// GetSingleBillingArtifact gets the Customer Statements with the specified CloudHealth Customer ID.
func (s *Client) GetSingleBillingArtifact(id int) (*BillingArtifact, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("customer_statements/%d?api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var billingArtifact = new(BillingArtifact)
	err = json.Unmarshal(responseBody, &billingArtifact)
	if err != nil {
		return nil, err
	}

	return billingArtifact, nil
}

// GetSingleBillingArtifacts gets all Customer Statements.
func (s *Client) GetSingleBillingArtifacts() (*BillingArtifacts, error) {
	billingArtifacts := new(BillingArtifacts)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"50"}, "api_key": {s.APIKey}}
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
		if len(ba.BillingArtifacts) < 50 {
			break
		}
		page++
	}
	return billingArtifacts, nil
}
