package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// AwsAccountAssignments represents all Assignments of AWS accounts enabled in CloudHealth with their details.
type AwsAccountAssignments struct {
	AwsAccountAssignments []AwsAccountAssignment `json:"aws_account_assignments"`
}

// AwsAccountAssignment represents the configuration of an AWS Account Assignment in CloudHealth with its details.
type AwsAccountAssignment struct {
	ID                  int    `json:"id"`
	OwnerID             string `json:"owner_id"`
	CustomerID          int    `json:"customer_id"`
	PayerAccountOwnerId string `json:"payer_account_owner_id"`
}

// GetSingleAwsAccountAssignment gets the details for the Assignment with specified ID.
func (s *Client) GetSingleAwsAccountAssignment(id int) (*AwsAccountAssignment, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_account_assignments/%d?api_key=%s", id, s.APIKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var awsaccountassignment = new(AwsAccountAssignment)
	err = json.Unmarshal(responseBody, &awsaccountassignment)
	if err != nil {
		return nil, err
	}

	return awsaccountassignment, nil
}

// GetAwsAccountAssignments gets all Assignments.
func (s *Client) GetAwsAccountAssignments() (*AwsAccountAssignments, error) {
	awsaccountassignments := new(AwsAccountAssignments)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"50"}, "api_key": {s.APIKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("aws_account_assignments/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var aa = new(AwsAccountAssignments)
		err = json.Unmarshal(responseBody, &aa)
		if err != nil {
			return nil, err
		}
		for _, a := range aa.AwsAccountAssignments {
			awsaccountassignments.AwsAccountAssignments = append(awsaccountassignments.AwsAccountAssignments, a)
		}
		if len(aa.AwsAccountAssignments) < 50 {
			break
		}
		page++
	}
	return awsaccountassignments, nil
}
