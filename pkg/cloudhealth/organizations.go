package cloudhealth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

// Organizations represents all Organizations enabled in CloudHealth with their configurations
type Organizations struct {
	Organizations []Organization `json:"organizations"`
}

// Organization represents the configuration of an Organization in CloudHealth
type Organization struct {
	ID                        string `json:"id"`
	ParentOrganizationID      string `json:"parent_organization_id"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	IdpName                   string `json:"idp_name"`
	FlexOrg                   bool   `json:"flex_org"`
	DefaultOrganization       bool   `json:"default_organization"`
	AssignedUsersCount        int    `json:"assigned_users_count"`
	NumAwsAccounts            int    `json:"num_aws_accounts"`
	NumAzureSubscriptions     int    `json:"num_azure_subscriptions"`
	NumGcpComputeProjects     int    `json:"num_gcp_compute_projects"`
	NumDataCenterAccounts     int    `json:"num_data_center_accounts"`
	NumVmwareCspOrganizations int    `json:"num_vmware_csp_organizations"`
}

// GetSingleOrganization gets the Organization with the specified
func (s *Client) GetSingleOrganization(id string) (*Organization, error) {
	// Set up the query parameters for the API
	params := url.Values{"org_id": {id}}

	relativeURL := fmt.Sprintf("v2/organizations?%s", params.Encode())

	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		return nil, err
	}

	var organizations Organizations
	err = json.Unmarshal(responseBody, &organizations)
	if err != nil {
		return nil, err
	}

	return &organizations.Organizations[0], nil
}

// GetOrganizations gets all Organizations listed in CloudHealth
func (s *Client) GetOrganizations() (*Organizations, error) {
	// Set variables we will need along the way
	var organizations Organizations
	var page, pageSize int = 1, 100

	// Loop for paging
	for {
		// Set up the query parameters for the API
		params := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(pageSize)}}

		// Set up the URL
		relativeURL := fmt.Sprintf("v2/organizations?%s", params.Encode())

		// Make the API call
		responseBody, err := getResponsePage(s, relativeURL)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response data into the Organizations struct
		err = json.Unmarshal(responseBody, &organizations)
		if err != nil {
			log.Fatalf("Error while putting JSON into a struct")
			return nil, err
		}

		// Check length of array in Organizations' struct to determine if we should break out of the loop
		if len(organizations.Organizations) < pageSize {
			break
		}

		// Increment page counter
		page++
	}

	return &organizations, nil
}
