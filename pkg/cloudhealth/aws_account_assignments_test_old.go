package cloudhealth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultAwsAccountAssignment = AwsAccountAssignment{
	ID:                  1234567890,
	OwnerID:             "9876543210",
	CustomerID:          1234,
	PayerAccountOwnerID: "9876543210",
}

var defaultAwsAccountAssignments = AwsAccountAssignments{
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

func TestGetSingleAwsAccountAssignment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("v1/aws_account_assignments/%d", defaultAwsAccountAssignment.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultAwsAccountAssignment)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedAwsAccountAssignment, err := c.GetSingleAwsAccountAssignment(1234567890)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if returnedAwsAccountAssignment.CustomerID != defaultAwsAccountAssignment.CustomerID {
		t.Errorf("GetAwsAccountAssignment() expected AwsAccountAssignmentID `%d`, got `%d`", defaultAwsAccountAssignment.CustomerID, returnedAwsAccountAssignment.CustomerID)
		return
	}
}

func TestGetAwsAccountAssignments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := "/aws_account_assignments/"
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(defaultAwsAccountAssignments)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedAwsAccountAssignments, err := c.GetAwsAccountAssignments()
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}
	if len(returnedAwsAccountAssignments.AwsAccountAssignments) != 2 {
		t.Errorf("All Aws Account Assignment statements have not been retrieved")
		return
	}
	if returnedAwsAccountAssignments.AwsAccountAssignments[0].CustomerID != defaultAwsAccountAssignments.AwsAccountAssignments[0].CustomerID && returnedAwsAccountAssignments.AwsAccountAssignments[1].CustomerID != defaultAwsAccountAssignments.AwsAccountAssignments[1].CustomerID {
		t.Errorf("Retrieved Aws Account Assignment statements don't match")
		return
	}
	return
}

func TestCreateAwsAccountAssignment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if r.URL.EscapedPath() != "/aws_account_assignments" {
			t.Errorf("Expected request to ‘/aws_account_assignments, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("Unable to read response body")
		}

		awsAccountAssignment := new(AwsAccountAssignment)
		err = json.Unmarshal(body, &awsAccountAssignment)
		if err != nil {
			t.Errorf("Unable to unmarshal awsAccountAssignment, got `%s`", body)
		}
		if awsAccountAssignment.CustomerID != 1234 {
			t.Errorf("Expected request to include awsAccountAssignment CustomerID ‘1234’, got ‘%v’", awsAccountAssignment.CustomerID)
		}
		awsAccountAssignment.ID = 1234567890
		js, err := json.Marshal(awsAccountAssignment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedAwsAccountAssignment, err := c.CreateAwsAccountAssignment(AwsAccountAssignment{
		CustomerID: 1234,
	})
	if err != nil {
		t.Errorf("CreateAwsAccountAssignment() returned an error: %s", err)
		return
	}
	if returnedAwsAccountAssignment.ID != 1234567890 {
		t.Errorf("CreateawsAccountAssignment() expected ID `1234567890`, got `%v`", returnedAwsAccountAssignment.ID)
		return
	}
}

func TestUpdateAwsAccountAssignment(t *testing.T) {
	updatedAwsAccountAssignment := defaultAwsAccountAssignment
	updatedAwsAccountAssignment.PayerAccountOwnerID = "4567"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/aws_account_assignments/%d", defaultAwsAccountAssignment.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(updatedAwsAccountAssignment)
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedAwsAccountAssignment, err := c.UpdateAwsAccountAssignment(updatedAwsAccountAssignment)
	if err != nil {
		t.Errorf("UpdateAwsAccountAssignment() returned an error: %s", err)
		return
	}
	if returnedAwsAccountAssignment.ID != updatedAwsAccountAssignment.ID {
		t.Errorf("UpdateAwsAccountAssignment() expected ID `%d`, got `%d`", defaultAwsAccountAssignment.ID, returnedAwsAccountAssignment.ID)
		return
	}
	if returnedAwsAccountAssignment.PayerAccountOwnerID == defaultAwsAccountAssignment.PayerAccountOwnerID {
		t.Errorf("UpdateAwsAccountAssignment() did not update the name")
		return
	}
}

func TestDeleteAwsAccountAssignment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/aws_account_assignments/%d", defaultAwsAccountAssignment.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
	}))
	defer ts.Close()

	c, err := NewClient("apiKey", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	err = c.DeleteAwsAccountAssignment(defaultAwsAccountAssignment.ID)
	if err != nil {
		t.Errorf("DeleteAwsAccountAssignment() returned an error: %s", err)
		return
	}
}
