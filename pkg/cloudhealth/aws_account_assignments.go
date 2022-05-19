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
	// Set up the URL
	relativeURL := fmt.Sprintf("v2/aws_account_assignments/%d", id)

	// Make the API call
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsAccountAssignment struct
	var awsaccountassignment AwsAccountAssignment
	err = json.Unmarshal(responseBody, &awsaccountassignment)
	if err != nil {
		return nil, err
	}

	return &awsaccountassignment, nil
}

// GetAwsAccountAssignments gets all Assignments.
func (s *Client) GetAwsAccountAssignments() (*AwsAccountAssignments, error) {
	// Set variables we will need along the way
	var awsaccountassignments AwsAccountAssignments
	var page, pageSize int = 1, 50

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v2/aws_account_assignments?%s", params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the AwsAccountAssignments struct
		err = json.Unmarshal(responseBody, &awsaccountassignments)
		if err != nil {
			return nil, err
		}

		// Check length of array in AwsAccountAssignments' struct to determine if we should break out of the loop
		if len(awsaccountassignments.AwsAccountAssignments) < pageSize {
			break
		}

		// Increment page counter
		page++
	}
	return &awsaccountassignments, nil
}

// CreateAwsAccountAssignment creates a new AwsAccountAssignment in CloudHealth.
func (s *Client) CreateAwsAccountAssignment(awsaccountassignment AwsAccountAssignment) (*AwsAccountAssignment, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v2/aws_account_assignments")

	// Make the API call
	responseBody, err := createResource(s, relativeURL, awsaccountassignment)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsAccountAssignment struct
	var returnedAwsAccountAssignment AwsAccountAssignment
	err = json.Unmarshal(responseBody, &returnedAwsAccountAssignment)
	if err != nil {
		return nil, err
	}

	return &returnedAwsAccountAssignment, nil
}

// UpdateAwsAccountAssignment updates an existing AwsAccountAssignment in CloudHealth.
func (s *Client) UpdateAwsAccountAssignment(awsaccountassignment AwsAccountAssignment) (*AwsAccountAssignment, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v2/aws_account_assignments/%d", awsaccountassignment.ID)

	// Make the API call
	responseBody, err := updateResource(s, relativeURL, awsaccountassignment)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsAccountAssignment struct
	var returnedAwsAccountAssignment AwsAccountAssignment
	err = json.Unmarshal(responseBody, &returnedAwsAccountAssignment)
	if err != nil {
		return nil, err
	}

	return &returnedAwsAccountAssignment, nil
}

// DeleteAwsAccountAssignment removes the AwsAccountAssignment with the specified CloudHealth ID.
func (s *Client) DeleteAwsAccountAssignment(id int) error {
	// Set up the URL
	relativeURL := fmt.Sprintf("v2/aws_account_assignments/%d", id)

	// Make the API call
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}
