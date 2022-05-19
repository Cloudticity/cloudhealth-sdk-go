package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// AwsAccounts represents all AWS Accounts enabled in CloudHealth with their configurations.
type AwsAccounts struct {
	AwsAccounts []AwsAccount `json:"aws_accounts"`
}

// AwsAccount represents the configuration of an AWS Account enabled in CloudHealth.
type AwsAccount struct {
	ID               int                      `json:"id,omitempty"`
	Name             string                   `json:"name"`
	OwnerID          string                   `json:"owner_id,omitempty"`
	HidePublicFields bool                     `json:"hide_public_fields,omitempty"`
	Region           string                   `json:"region,omitempty"`
	CreatedAt        time.Time                `json:"created_at,omitempty"`
	UpdatedAt        time.Time                `json:"updated_at,omitempty"`
	AccountType      string                   `json:"account_type,omitempty"`
	VpcOnly          bool                     `json:"vpc_only,omitempty"`
	ClusterName      string                   `json:"cluster_name,omitempty"`
	Status           AwsAccountStatus         `json:"status,omitempty"`
	Authentication   AwsAccountAuthentication `json:"authentication,omitempty"`
	Tags             []map[string]string      `json:"tags,omitempty"`
}

// AwsAccountStatus represents the status details for AWS integration.
type AwsAccountStatus struct {
	Level      string    `json:"level"`
	LastUpdate time.Time `json:"last_update,omitempty"`
}

// AwsAccountAuthentication represents the authentication details for AWS integration.
type AwsAccountAuthentication struct {
	Protocol             string `json:"protocol,omitempty"`
	AccessKey            string `json:"access_key,omitempty"`
	SecretKey            string `json:"secret_key,omitempty"`
	AssumeRoleArn        string `json:"assume_role_arn,omitempty"`
	AssumeRoleExternalID string `json:"assume_role_external_id,omitempty"`
}

// AwsExternalID is used to enable integration with AWS via IAM Roles.
type AwsExternalID struct {
	ExternalID string `json:"generated_external_id,omitempty"`
}

// GetSingleAwsAccount gets the AWS Account with the specified CloudHealth Account ID.
func (s *Client) GetSingleAwsAccount(id int) (*AwsAccount, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/aws_accounts/%d", id)

	// Make the API call
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsAccount struct
	var account AwsAccount
	err = json.Unmarshal(responseBody, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAwsAccounts gets all AWS Accounts enabled in CloudHealth.
func (s *Client) GetAwsAccounts() (*AwsAccounts, error) {
	// Set variables we will need along the way
	var awsaccounts AwsAccounts
	var page, pageSize int = 1, 100

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v1/aws_accounts?%s", params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the AWSAccounts struct
		err = json.Unmarshal(responseBody, &awsaccounts)
		if err != nil {
			return nil, err
		}

		// Check length of array in AwsAccounts' struct to determine if we should break out of the loop
		if len(awsaccounts.AwsAccounts) < pageSize {
			break
		}

		// Increment page counter
		page++
	}
	return &awsaccounts, nil
}

// CreateAwsAccount enables a new AWS Account in CloudHealth.
func (s *Client) CreateAwsAccount(account AwsAccount) (*AwsAccount, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/aws_accounts")

	// Make the API call
	responseBody, err := createResource(s, relativeURL, account)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsAccount struct
	var returnedAccount AwsAccount
	err = json.Unmarshal(responseBody, &returnedAccount)
	if err != nil {
		return nil, err
	}

	return &returnedAccount, nil
}

// UpdateAwsAccount updates an existing AWS Account in CloudHealth.
func (s *Client) UpdateAwsAccount(account AwsAccount) (*AwsAccount, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/aws_accounts/%d", account.ID)

	// Make the API call
	responseBody, err := updateResource(s, relativeURL, account)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsAccount struct
	var returnedAccount AwsAccount
	err = json.Unmarshal(responseBody, &returnedAccount)
	if err != nil {
		return nil, err
	}

	return &returnedAccount, nil
}

// DeleteAwsAccount removes the AWS Account with the specified CloudHealth ID.
func (s *Client) DeleteAwsAccount(id int) error {
	// Set up the URL
	relativeURL := fmt.Sprintf("aws_accounts/%d", id)

	// Make the API call
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}

// GetAwsExternalID gets the AWS External ID tied to the CloudHealth Account.
func (s *Client) GetAwsExternalID(id int) (*AwsExternalID, error) {
	// Set up the URL
	relativeURL := fmt.Sprintf("v1/aws_accounts/%d/generate_external_id", id)

	// Make the API call
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into the AwsExternalId struct
	var externalId AwsExternalID
	err = json.Unmarshal(responseBody, &externalId)
	if err != nil {
		return nil, err
	}

	return &externalId, nil
}
