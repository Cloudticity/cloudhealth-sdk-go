package cloudhealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

// ErrAwsAccountNotFound is returned when an AWS Account doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if exists).
var ErrAwsAccountNotFound = errors.New("AWS Account not found")

// GetAwsAccount gets the AWS Account with the specified CloudHealth Account ID. (deprecated, will be removed in future, kept only to not break anything)
func (s *Client) GetAwsAccount(id int) (*AwsAccount, error) {
	return s.GetSingleAwsAccount(id)
}

// GetSingleAwsAccount gets the AWS Account with the specified CloudHealth Account ID.
func (s *Client) GetSingleAwsAccount(id int) (*AwsAccount, error) {
	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", id, s.ApiKey))

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrAwsAccountNotFound
		}
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

	body, _ := json.Marshal(account)

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts?api_key=%s", s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		var account = new(AwsAccount)
		err = json.Unmarshal(responseBody, &account)
		if err != nil {
			return nil, err
		}

		return account, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("Bad Request. Please check if a AWS Account with this name `%s` already exists", account.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// UpdateAwsAccount updates an existing AWS Account in CloudHealth.
func (s *Client) UpdateAwsAccount(account AwsAccount) (*AwsAccount, error) {

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", account.ID, s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	body, _ := json.Marshal(account)

	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer((body)))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var account = new(AwsAccount)
		err = json.Unmarshal(responseBody, &account)
		if err != nil {
			return nil, err
		}

		return account, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusNotFound:
		return nil, ErrAwsAccountNotFound
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("Bad Request. Please check if a AWS Account with this name `%s` already exists", account.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// DeleteAwsAccount removes the AWS Account with the specified CloudHealth ID.
func (s *Client) DeleteAwsAccount(id int) error {

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", id, s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("DELETE", url.String(), nil)

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrAwsAccountNotFound
	case http.StatusUnauthorized:
		return ErrClientAuthenticationError
	default:
		return fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
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
