package cloudhealth

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AwsExternalID is used to enable integration with AWS via IAM Roles.
type AwsExternalID struct {
	ExternalID string `json:"generated_external_id"`
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
