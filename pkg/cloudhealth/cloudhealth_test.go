package cloudhealth

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

// We'll be able to store suite-wide variables and add methods to this test suite struct
type CloudHealthTestSuite struct {
	suite.Suite
	ts                           *httptest.Server
	defaultAwsAccountAssignment  AwsAccountAssignment
	defaultAwsAccountAssignments AwsAccountAssignments
	createTestServer             func(responseBody interface{}) *httptest.Server
}

// This will run right before the test starts and receives the suite and test names as input
func (suite *CloudHealthTestSuite) BeforeTest(suiteName, testName string) {}

// This will run after test finishes and receives the suite and test names as input
func (suite *CloudHealthTestSuite) AfterTest(suiteName, testName string) {}

// This will run before the tests in the suite are run
func (suite *CloudHealthTestSuite) SetupSuite() {}

// This will run before each test in the suite
func (suite *CloudHealthTestSuite) SetupTest() {
	suite.defaultAwsAccountAssignment = AwsAccountAssignment{
		ID:                  1234567890,
		OwnerID:             "9876543210",
		CustomerID:          1234,
		PayerAccountOwnerID: "9876543210",
	}

	suite.defaultAwsAccountAssignments = AwsAccountAssignments{
		AwsAccountAssignments: []AwsAccountAssignment{
			{
				ID:                  1234567890,
				OwnerID:             "9876543210",
				CustomerID:          1234,
				PayerAccountOwnerID: "9876543210",
			},
			{
				ID:                  1234567890,
				OwnerID:             "9876543210",
				CustomerID:          1234,
				PayerAccountOwnerID: "9876543210",
			},
		},
	}

	suite.createTestServer = func(responseBody interface{}) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set HTTP 200 for response header
			w.WriteHeader(http.StatusOK)

			// Set response body
			body, _ := json.Marshal(responseBody)
			w.Write(body)
		}))
	}
}

// This is an example test that will always succeed
func (suite *CloudHealthTestSuite) TestExample() {
	suite.Equal(true, true)
}

// We need this function to kick off the test suite, otherwise
// "go test" won't know about our tests
func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(CloudHealthTestSuite))
}
