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
	relativeURL, _ := url.Parse(fmt.Sprintf("price_book_account_assignments/%d?api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var accountPriceBookAssignment = new(AccountPriceBookAssignment)
	err = json.Unmarshal(responseBody, &accountPriceBookAssignment)
	if err != nil {
		return nil, err
	}

	return accountPriceBookAssignment, nil
}

// GetAccountPriceBookAssignments gets all Assignments.
func (s *Client) GetAccountPriceBookAssignments() (*AccountPriceBookAssignments, error) {
	accountPriceBookAssignments := new(AccountPriceBookAssignments)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"50"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("price_book_account_assignments/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var ac = new(AccountPriceBookAssignments)
		err = json.Unmarshal(responseBody, &ac)
		if err != nil {
			return nil, err
		}
		for _, a := range ac.AccountPriceBookAssignments {
			accountPriceBookAssignments.AccountPriceBookAssignments = append(accountPriceBookAssignments.AccountPriceBookAssignments, a)
		}
		if len(ac.AccountPriceBookAssignments) < 50 {
			break
		}
		page++
	}
	return accountPriceBookAssignments, nil
}
