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
	ID               int                      `json:"id"`
	Name             string                   `json:"name"`
	OwnerID          string                   `json:"owner_id"`
	HidePublicFields bool                     `json:"hide_public_fields"`
	Region           string                   `json:"region"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
	AccountType      string                   `json:"account_type"`
	VpcOnly          bool                     `json:"vpc_only"`
	ClusterName      string                   `json:"cluster_name"`
	Status           AwsAccountStatus         `json:"status"`
	Authentication   AwsAccountAuthentication `json:"authentication"`
}

// AwsAccountStatus represents the status details for AWS integration.
type AwsAccountStatus struct {
	Level      string    `json:"level"`
	LastUpdate time.Time `json:"last_update"`
}

// AwsAccountAuthentication represents the authentication details for AWS integration.
type AwsAccountAuthentication struct {
	Protocol             string `json:"protocol"`
	AccessKey            string `json:"access_key"`
	SecreyKey            string `json:"secret_key"`
	AssumeRoleArn        string `json:"assume_role_arn"`
	AssumeRoleExternalID string `json:"assume_role_external_id"`
}

// AwsExternalID is used to enable integration with AWS via IAM Roles.
type AwsExternalID struct {
	ExternalID string `json:"generated_external_id"`
}

// GetAwsAccount gets the AWS Account with the specified CloudHealth Account ID. (deprecated, will be removed in future, kept only to not break anything)
func (s *Client) GetAwsAccount(id int) (*AwsAccount, error) {
	return s.GetSingleAwsAccount(id)
}

// GetSingleAwsAccount gets the AWS Account with the specified CloudHealth Account ID.
func (s *Client) GetSingleAwsAccount(id int) (*AwsAccount, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", id, s.ApiKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var account = new(AwsAccount)
	err = json.Unmarshal(responseBody, &account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAwsAccounts gets all AWS Accounts enabled in CloudHealth.
func (s *Client) GetAwsAccounts() (*AwsAccounts, error) {
	awsaccounts := new(AwsAccounts)
	page := 1
	for {
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {"100"}, "api_key": {s.ApiKey}}
		relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/?%s", params.Encode()))
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}
		var acts = new(AwsAccounts)
		err = json.Unmarshal(responseBody, &acts)
		if err != nil {
			return nil, err
		}
		for _, p := range acts.AwsAccounts {
			awsaccounts.AwsAccounts = append(awsaccounts.AwsAccounts, p)
		}
		if len(acts.AwsAccounts) < 100 {
			break
		}
		page++
	}
	return awsaccounts, nil
}

// CreateAwsAccount enables a new AWS Account in CloudHealth.
func (s *Client) CreateAwsAccount(account AwsAccount) (*AwsAccount, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts?api_key=%s", s.ApiKey))

	responseBody, err := createResource(s, relativeURL, account)
	if err != nil {
		return nil, err
	}

	var returnedAccount = new(AwsAccount)
	err = json.Unmarshal(responseBody, &returnedAccount)
	if err != nil {
		return nil, err
	}

	return returnedAccount, nil
}

// UpdateAwsAccount updates an existing AWS Account in CloudHealth.
func (s *Client) UpdateAwsAccount(account AwsAccount) (*AwsAccount, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", account.ID, s.ApiKey))

	responseBody, err := updateResource(s, relativeURL, account)
	if err != nil {
		return nil, err
	}

	var returnedAccount = new(AwsAccount)
	err = json.Unmarshal(responseBody, &returnedAccount)
	if err != nil {
		return nil, err
	}

	return returnedAccount, nil
}

// DeleteAwsAccount removes the AWS Account with the specified CloudHealth ID.
func (s *Client) DeleteAwsAccount(id int) error {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", id, s.ApiKey))
	_, err := deleteResource(s, relativeURL)
	if err != nil {
		return err
	}

	return nil
}

// GetAwsExternalID gets the AWS External ID tied to the CloudHealth Account.
func (s *Client) GetAwsExternalID(id int) (*AwsExternalID, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d/generate_external_id?api_key=%s", id, s.ApiKey))
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var extid = new(AwsExternalID)
	err = json.Unmarshal(responseBody, &extid)
	if err != nil {
		return nil, err
	}

	return extid, nil
}
