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
	ID                  int    `json:"id,omitempty"`
	OwnerID             string `json:"owner_id"`
	CustomerID          int    `json:"customer_id"`
	PayerAccountOwnerID string `json:"payer_account_owner_id,omitempty"`
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

// CreateAwsAccountAssignment creates a new AwsAccountAssignment in CloudHealth.
func (s *Client) CreateAwsAccountAssignment(awsaccountassignment AwsAccountAssignment) (*AwsAccountAssignment, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_account_assignments?api_key=%s", s.APIKey))

	responseBody, err := createResource(s, relativeURL, awsaccountassignment)
	if err != nil {
		return nil, err
	}

	var returnedAwsAccountAssignment = new(AwsAccountAssignment)
	err = json.Unmarshal(responseBody, &returnedAwsAccountAssignment)
	if err != nil {
		return nil, err
	}

	return returnedAwsAccountAssignment, nil
}

// UpdateAwsAccountAssignment updates an existing AwsAccountAssignment in CloudHealth.
func (s *Client) UpdateAwsAccountAssignment(awsaccountassignment AwsAccountAssignment) (*AwsAccountAssignment, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_account_assignments/%d?api_key=%s", awsaccountassignment.ID, s.APIKey))

	responseBody, err := updateResource(s, relativeURL, awsaccountassignment)
	if err != nil {
		return nil, err
	}

	var returnedAwsAccountAssignment = new(AwsAccountAssignment)
	err = json.Unmarshal(responseBody, &returnedAwsAccountAssignment)
	if err != nil {
		return nil, err
	}

	return returnedAwsAccountAssignment, nil
}

// DeleteAwsAccountAssignment removes the AwsAccountAssignment with the specified CloudHealth ID.
func (s *Client) DeleteAwsAccountAssignment(id int) error {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_account_assignments/%d?api_key=%s", id, s.APIKey))
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}
