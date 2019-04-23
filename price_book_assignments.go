package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// CustomerPriceBookAssignments represents all assignments to a Custom Price Book for all Customers.
type CustomerPriceBookAssignments struct {
	CustomerPriceBookAssignments []CustomerPriceBookAssignment `json:"price_book_assignments"`
}

// CustomerPriceBookAssignment represents the configuration of a Customer Price Book assignment to a Customer.
type CustomerPriceBookAssignment struct {
	ID                int       `json:"id"`
	PriceBookID       int       `json:"price_book_id"`
	TargetClientAPIID int       `json:"target_client_api_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// GetSingleCustomerPriceBookAssignment gets the details for the Assignment with specified ID.
func (s *Client) GetSingleCustomerPriceBookAssignment(id int) (*CustomerPriceBookAssignment, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("price_book_assignments/%d?api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var customerPriceBookAssignment = new(CustomerPriceBookAssignment)
	err = json.Unmarshal(responseBody, &customerPriceBookAssignment)
	if err != nil {
		return nil, err
	}

	return customerPriceBookAssignment, nil
}

// GetCustomerPriceBookAssignments gets all Assignments.
func (s *Client) GetCustomerPriceBookAssignments() (*CustomerPriceBookAssignments, error) {
	customerPriceBookAssignments := new(CustomerPriceBookAssignments)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"50"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("price_book_assignments/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var cp = new(CustomerPriceBookAssignments)
		err = json.Unmarshal(responseBody, &cp)
		if err != nil {
			return nil, err
		}
		for _, c := range cp.CustomerPriceBookAssignments {
			customerPriceBookAssignments.CustomerPriceBookAssignments = append(customerPriceBookAssignments.CustomerPriceBookAssignments, c)
		}
		if len(cp.CustomerPriceBookAssignments) < 50 {
			break
		}
		page++
	}
	return customerPriceBookAssignments, nil
}

// DeleteCustomerPriceBookAssignment removes the Customer Price Book Assignment with the specified CloudHealth ID.
func (s *Client) DeleteCustomerPriceBookAssignment(id int) error {
	relativeURL, _ := url.Parse(fmt.Sprintf("price_book_assignments/%d?api_key=%s", id, s.APIKey))
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}
