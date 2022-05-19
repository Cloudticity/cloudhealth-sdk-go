package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Customers represents all Customers enabled in CloudHealth with their configurations.
type Customers struct {
	Customers []Customer `json:"customers"`
}

// Customer represents the configuration of a Customer in CloudHealth.
type Customer struct {
	ID                          int                                 `json:"id,omitempty"`
	Name                        string                              `json:"name"`
	Classification              string                              `json:"classification,omitempty"`
	MarginPercentage            float64                             `json:"margin_percentage,omitempty"`
	CreatedAt                   time.Time                           `json:"created_at,omitempty"`
	UpdatedAt                   time.Time                           `json:"updated_at,omitempty"`
	GeneralExternalID           string                              `json:"generated_external_id,omitempty"`
	PartnerBillingConfiguration CustomerPartnerBillingConfiguration `json:"partner_billing_configuration,omitempty"`
	Address                     CustomerAddress                     `json:"address"`
	BillingConfiguration        CustomerBillingConfiguration        `json:"billing_configuration,omitempty"`
	Tags                        []map[string]string                 `json:"tags,omitempty"`
}

// CustomerPartnerBillingConfiguration represents partner billing details of a Customer.
type CustomerPartnerBillingConfiguration struct {
	Enabled bool   `json:"enabled"`
	Folder  string `json:"folder,omitempty"`
}

// CustomerBillingConfiguration represents billing details of a Customer.
type CustomerBillingConfiguration struct {
	Status string `json:"status"`
}

// CustomerAddress represents address details of a Customer.
type CustomerAddress struct {
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipcode"`
	Country string `json:"country"`
}

// GetSingleCustomer gets the Customer with the specified CloudHealth Customer ID.
func (s *Client) GetSingleCustomer(id int) (*Customer, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/customers/%d", id)

	// Make the API call
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the Customer struct
	var customer Customer
	err = json.Unmarshal(responseBody, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

// GetCustomers gets all AWS Accounts enabled in CloudHealth.
func (s *Client) GetCustomers() (*Customers, error) {
	// Set variables we will need along the way
	var customers Customers
	var page, pageSize int = 1, 100

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v1/customers?%s", params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the Customers struct
		err = json.Unmarshal(responseBody, &customers)
		if err != nil {
			return nil, err
		}

		// Check length of array in Customers' struct to determine if we should break out of the loop
		if len(customers.Customers) < pageSize {
			break
		}

		// Increment page counter
		page++
	}
	return &customers, nil
}

// CreateCustomer creates a new Customer in CloudHealth.
func (s *Client) CreateCustomer(customer Customer) (*Customer, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/customers")

	// Make the API call
	responseBody, err := createResource(s, relativeURL, customer)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the Organizations struct
	var returnedCustomer Customer
	err = json.Unmarshal(responseBody, &returnedCustomer)
	if err != nil {
		return nil, err
	}

	return &returnedCustomer, nil
}

// UpdateCustomer updates an existing Customer in CloudHealth.
func (s *Client) UpdateCustomer(customer Customer) (*Customer, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/customers/%d", customer.ID)

	// Make the API call
	responseBody, err := updateResource(s, relativeURL, customer)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the Organizations struct
	var returnedcustomer Customer
	err = json.Unmarshal(responseBody, &returnedcustomer)
	if err != nil {
		return nil, err
	}

	return &returnedcustomer, nil
}

// DeleteCustomer removes the Customer with the specified CloudHealth ID.
func (s *Client) DeleteCustomer(id int) error {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/customers/%d", id)

	// Make the API call
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}
