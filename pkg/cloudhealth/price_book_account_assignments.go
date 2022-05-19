package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// AccountPriceBookAssignments represents all assignments to a Custom Price Book for all Accounts.
type AccountPriceBookAssignments struct {
	AccountPriceBookAssignments []AccountPriceBookAssignment `json:"price_book_account_assignments"`
}

// AccountPriceBookAssignment represents the configuration of a Customer Price Book assignment to a Customer.
type AccountPriceBookAssignment struct {
	ID                    int         `json:"id"`
	TargetClientAPIID     int         `json:"target_client_api_id"`
	PriceBookAssignmentID int         `json:"price_book_assignment_id"`
	BillingAccountOwnerID interface{} `json:"billing_account_owner_id"`
}

// GetSingleAccountPriceBookAssignment gets the details for the Assignment with specified ID.
func (s *Client) GetSingleAccountPriceBookAssignment(id int) (*AccountPriceBookAssignment, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/price_book_account_assignments/%d", id)

	// Make the API call
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AccountPriceBookAssignment struct
	var accountPriceBookAssignment AccountPriceBookAssignment
	err = json.Unmarshal(responseBody, &accountPriceBookAssignment)
	if err != nil {
		return nil, err
	}

	return &accountPriceBookAssignment, nil
}

// GetAccountPriceBookAssignments gets all Assignments.
func (s *Client) GetAccountPriceBookAssignments() (*AccountPriceBookAssignments, error) {
	// Set variables we will need along the way
	var accountPriceBookAssignments AccountPriceBookAssignments
	var page, pageSize int = 1, 50

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v1/price_book_account_assignments?%s", params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the AccountPriceBookAssignments struct
		err = json.Unmarshal(responseBody, &accountPriceBookAssignments)
		if err != nil {
			return nil, err
		}

		// Check length of array in AccountPriceBookAssignments' struct to determine if we should break out of the loop
		if len(accountPriceBookAssignments.AccountPriceBookAssignments) < pageSize {
			break
		}

		// Increment page counter
		page++
	}
	return &accountPriceBookAssignments, nil
}
