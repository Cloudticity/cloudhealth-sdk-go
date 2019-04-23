[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/nextgenhealthcare/cloudhealth-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nextgenhealthcare/cloudhealth-sdk-go)](https://goreportcard.com/report/github.com/nextgenhealthcare/cloudhealth-sdk-go)

# CloudHealth API in Go

A Go wrapper for the [CloudHealth Cloud Management Platform](https://www.cloudhealthtech.com/) API.

## Getting Started

Run the following command to retrieve the SDK:

```bash
go get -u github.com/nextgenhealthcare/cloudhealth-sdk-go
```

You will also need an API Key from CloudHealth. For more information, see [Getting Your API Key](http://apidocs.cloudhealthtech.com/#documentation_getting-your-api-key)

```go
import cloudhealth "github.com/nextgenhealthcare/cloudhealth-sdk-go"

client, _ := cloudhealth.NewClient("api_key", "https://chapi.cloudhealthtech.com/v1/")

account, err := client.GetSingleAwsAccount(1234567890)
if err == cloudhealth.ErrAwsAccountNotFound {
	log.Fatalf("AWS Account not found: %s\n", err)
}
if err != nil {
	log.Fatalf("Unknown error: %s\n", err)
}

log.Printf("AWS Account %s\n", account.Name)
```

## Available Endpoints

| Endpoint | HTTP Method | SDK Method | Description | Status |
| -- | -- | -- | -- | -- |
| `/aws_accounts` | `POST` | `CreateAwsAccount()` | Enable AWS Account | :heavy_check_mark: |
| `/aws_accounts` | `GET` | `GetAwsAccounts()` | AWS Accounts in CloudHealth  | :heavy_check_mark: |
| `/aws_accounts/:id` | `GET` | `GetSingleAwsAccount()` | Single AWS Account | :heavy_check_mark: |
| `/aws_accounts/:id` | `PUT` | `UpdateAwsAccount()` | Update Existing AWS Account | :heavy_check_mark: |
| `/aws_accounts/:id` | `DELETE` | `DeleteAwsAccount()` | Delete AWS Account | :heavy_check_mark: |
| `/aws_accounts/:id/generate_external_id` | `GET` | `GetAwsExternalID()` | Get External ID | :heavy_check_mark: |
| `/customers` | `POST` | `CreateCustomer()` | Create Partner Customer  | :heavy_check_mark: |
| `/customers/:id` | `PUT` | `UpdateCustomer()` | Modify Existing Customer | :heavy_check_mark:  |
| `/customers/:id` | `DELETE` | `DeleteCustomer` | Delete Existing Customer  | :heavy_check_mark: |
| `/customers/:id` | `GET` | `GetSingleCustomer()` | Get Single Customer | :heavy_check_mark: |
| `/customers` | `GET` | `GetCustomers` | Get All Customers | :heavy_check_mark: |
| `/customer_statements` | `GET` | `GetSingleBillingArtifacts()` | Statement for Single Customer | :heavy_check_mark: |
| `/customer_statements` | `GET` | `GetBillingArtifacts()` | Statements for All Customers | :heavy_check_mark: |
| `/aws_account_assignments` | `POST` | | Create AWS Account Assignment |  |
| `/aws_account_assignments` | `GET` | | Read All AWS Account Assignments | :heavy_check_mark: |
| `/aws_account_assignments/:id` | `GET` | | Read Single AWS Account Assignment | :heavy_check_mark: |
| `/aws_account_assignments/:id` | `PUT` | | Update AWS Account Assignment | :heavy_check_mark:  |
| `/aws_account_assignments/:id` | `DELETE` | | Delete AWS Account Assignment | :heavy_check_mark: |
| `/price_book_assignments` | `GET` | | Read qll Customer Price Book Assignments | :heavy_check_mark: |
| `/price_book_assignments/:id` | `GET` | | Read Single Customer Price Book Assignment | :heavy_check_mark: |
| `/price_book_assignments/:id` | `DELETE` | | Delete Customer Price Book Assignment | :heavy_check_mark: |

## Contributing

Any and all contributions are welcome. Please don't hesitate to submit an issue or pull request.

## Roadmap

The initial release was focused on being consumed by a Terraform provider in AWS environments such as support for managing AWS Accounts in CloudHealth.

The current state is a global SDK with all major primitives to set up your CloudHealth environment.

## Development

Build and Install with `make` or `make build`.
Run `gofmt` with `make fmt`.
Run unit tests with `make test`.
Run coverage tests with `make cover`.
