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
	relativeURL, _ := url.Parse(fmt.Sprintf("customers/%d?api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var customer = new(Customer)
	err = json.Unmarshal(responseBody, &customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

// GetCustomers gets all AWS Accounts enabled in CloudHealth.
func (s *Client) GetCustomers() (*Customers, error) {
	customers := new(Customers)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"50"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("customers/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var csts = new(Customers)
		err = json.Unmarshal(responseBody, &csts)
		if err != nil {
			return nil, err
		}
		for _, a := range csts.Customers {
			customers.Customers = append(customers.Customers, a)
		}
		if len(csts.Customers) < 50 {
			break
		}
		page++
	}
	return customers, nil
}

// CreateCustomer creates a new Customer in CloudHealth.
func (s *Client) CreateCustomer(customer Customer) (*Customer, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("customers?api_key=%s", s.APIKey))

	responseBody, err := createResource(s, relativeURL, customer)
	if err != nil {
		return nil, err
	}

	var returnedCustomer = new(Customer)
	err = json.Unmarshal(responseBody, &returnedCustomer)
	if err != nil {
		return nil, err
	}

	return returnedCustomer, nil
}

// UpdateCustomer updates an existing Customer in CloudHealth.
func (s *Client) UpdateCustomer(customer Customer) (*Customer, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("customers/%d?api_key=%s", customer.ID, s.APIKey))

	responseBody, err := updateResource(s, relativeURL, customer)
	if err != nil {
		return nil, err
	}

	var returnedcustomer = new(Customer)
	err = json.Unmarshal(responseBody, &returnedcustomer)
	if err != nil {
		return nil, err
	}

	return returnedcustomer, nil
}

// DeleteCustomer removes the Customer with the specified CloudHealth ID.
func (s *Client) DeleteCustomer(id int) error {
	relativeURL, _ := url.Parse(fmt.Sprintf("customers/%d?api_key=%s", id, s.APIKey))
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}
