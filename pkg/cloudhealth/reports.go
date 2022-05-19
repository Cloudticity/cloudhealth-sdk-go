package cloudhealth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// AWSCostHistoryRequestOptions represents the possible options to specify when making a request against the cost history report
type AWSCostHistoryRequestOptions struct {
	interval           string
	measures           string
	clientApiId        string
	selectedDimensions string
	rejectedDimensions string
	targetAWSAccountId string
	time               string
}

// AWSCostHistoryReport represents the details of a Cost History Report for the AWS Service Category in CloudHealth
type AWSCostHistoryReport struct {
	BillDropInfo         []interface{}                    `json:"bill_drop_info"`
	CubeID               string                           `json:"cube_id"`
	Data                 [][]float64                      `json:"data"`
	Dimensions           []AWSCostHistoryReportDimensions `json:"dimensions"`
	EnableDpPopover      bool                             `json:"enable_dp_popover"`
	Filters              []string                         `json:"filters"`
	Interval             string                           `json:"interval"`
	Measures             []AWSCostHistoryReportMeasures   `json:"measures"`
	Report               string                           `json:"report"`
	Status               string                           `json:"status"`
	UpdatedAt            time.Time                        `json:"updated_at,omitempty"`
	VisualizationOptions interface{}                      `json:"visualization_options"`
}

// AWSCostHistoryReportDimensions are the available dimensions of the Cost History Report with the AWS Service Category in CloudHealth
type AWSCostHistoryReportDimensions struct {
	AwsServiceCategory []AWSCostHistoryReportAwsServiceCategory `json:"AWS-Service-Category"`
}

// AWSCostHistoryReportAwsServiceCategory represents category details that pertain to the AWS Service Category in CloudHealth
type AWSCostHistoryReportAwsServiceCategory struct {
	Direct    bool        `json:"direct"`
	Excluded  interface{} `json:"excluded"`
	Extended  bool        `json:"extended"`
	Label     string      `json:"label"`
	Name      string      `json:"name"`
	Parent    int64       `json:"parent"`
	Populated interface{} `json:"populated"`
	SortOrder interface{} `json:"sort_order"`
}

// AWSCostHistoryReportMeasures represents the measures in the report of the AWS Service Category in CloudHealth
type AWSCostHistoryReportMeasures struct {
	Label    string                               `json:"label"`
	Metadata AWSCostHistoryReportMeasuresMetadata `json:"metadata"`
	Name     string                               `json:"name"`
}

// AWSCostHistoryReportMeasuresMetadata represents the available measures' metadata of the AWS Service Category in CloudHealth
type AWSCostHistoryReportMeasuresMetadata struct {
	AncillaryCaches   []string `json:"ancillary_caches"`
	Label             string   `json:"label"`
	SupportsDrilldown bool     `json:"supports_drilldown"`
	Type              string   `json:"type"`
	Units             string   `json:"units"`
}

// GetAWSCostHistoryReport gets the Cost History Report
func (s *Client) GetAWSCostHistoryReport(requestOptions *AWSCostHistoryRequestOptions) (*AWSCostHistoryReport, error) {
	// Set up the base URL
	relativeURL := fmt.Sprintf("olap_reports/cost/history?dimensions[]=AWS-Service-Category")

	// Parse request options
	if requestOptions.measures != "" {
		relativeURL = fmt.Sprintf("%s&measures[]=%s", relativeURL, requestOptions.measures)
	} else {
		return nil, errors.New("the `measures` property is required and cannot be blank")
	}

	if requestOptions.interval != "" {
		relativeURL = fmt.Sprintf("%s&interval=%s", relativeURL, requestOptions.interval)
	} else {
		return nil, errors.New("the `interval` property is required and cannot be blank")
	}

	if requestOptions.time != "" {
		relativeURL = fmt.Sprintf("%s&filters[]=time:select:%s", relativeURL, requestOptions.time)
	} else {
		return nil, errors.New("the `time` property is required and cannot be blank")
	}

	if requestOptions.clientApiId != "" {
		relativeURL = fmt.Sprintf("%s&client_api_id=%s", relativeURL, requestOptions.clientApiId)
	}

	if requestOptions.selectedDimensions != "" {
		relativeURL = fmt.Sprintf("%s&filters[]=AWS-Service-Category:select:%s", relativeURL, requestOptions.selectedDimensions)
	}

	if requestOptions.rejectedDimensions != "" {
		relativeURL = fmt.Sprintf("%s&filters[]=AWS-Service-Category:reject:%s", relativeURL, requestOptions.rejectedDimensions)
	}

	if requestOptions.targetAWSAccountId != "" {
		relativeURL = fmt.Sprintf("%s&filters[]=AWS-Account:select:%s", relativeURL, requestOptions.targetAWSAccountId)
	}

	// Make the API call
	responseBody, err := getResponsePage(s, relativeURL)
	if err != nil {
		fmt.Println("Error while calling CloudHealth API")
		return nil, err
	}

	// Unmarshal the response data into the Customer struct
	var costReport AWSCostHistoryReport
	err = json.Unmarshal(responseBody, &costReport)
	if err != nil {
		return nil, err
	}

	return &costReport, nil
}
