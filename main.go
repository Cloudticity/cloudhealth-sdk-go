package main

import (
	"fmt"
	"github.com/Cloudticity/cloudhealth-sdk-go/pkg/cloudhealth"
	"log"
)

func main() {
	c, err := cloudhealth.NewClient("Bearer 22fb4d3c-ada4-47da-84ae-2214379bed9d", "https://chapi.cloudhealthtech.com/")
	if err != nil {
		log.Fatalf("Got error setting up new client: %s", err)
	}

	report, err := c.GetAwsAccounts()
	if err != nil {
		log.Fatalf("Got error requesting report: %s", err)
	}

	fmt.Printf("%+v\n", report)
}
