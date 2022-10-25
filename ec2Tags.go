//
// Copyright (c) 2021-2022 Tenebris Technologies Inc.
//

package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"log2sqs/config"
)

func ec2Tags() {
	var awsCredentials *credentials.Credentials
	var awsConfig *aws.Config
	var instanceID = ""

	log.Printf("Reading AWS EC2 instance metadata...")

	// Get instance ID
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.Printf("Error retriving EC2 instance ID")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading HTTP response for EC2 instance ID")
		return
	} else {
		instanceID = string(body)
		addFields["_ec2_instanceID"] = instanceID
	}
	_ = resp.Body.Close()

	// Get hostname
	resp, err = http.Get("http://169.254.169.254/latest/meta-data/hostname")
	if err != nil {
		log.Printf("Error retriving EC2 hostname")
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading HTTP response for EC2 hostname")
		return
	} else {
		addFields["_ec2_hostname"] = string(body)
	}
	_ = resp.Body.Close()

	// Set up credentials for AWS API
	if config.AWSID == "role" {
		// Assume EC2 instance with permissions granted through IAM
		awsConfig = &aws.Config{
			Region: aws.String(config.AWSRegion),
		}
	} else {
		// Use credentials from configuration
		awsCredentials = credentials.NewStaticCredentials(config.AWSID, config.AWSKey, "")
		awsConfig = &aws.Config{
			Region:      aws.String(config.AWSRegion),
			Credentials: awsCredentials,
		}
	}

	// Initialize a session
	awsSession := session.Must(session.NewSession(awsConfig))

	// Create new EC2 client
	ec2Svc := ec2.New(awsSession)

	// Configure Input
	input := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	}

	ec2Info, err := ec2Svc.DescribeTags(input)
	if err != nil {
		log.Printf("Error retrieving EC2 addFields: %s", err.Error())
		return
	}

	// Iterate over EC2 tags and add to addFields[]
	for _, tag := range ec2Info.Tags {
		if tag.Key != nil && tag.Value != nil {
			addFields["_ec2_tag_"+*tag.Key] = *tag.Value
		}
	}
}
