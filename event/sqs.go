//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package event

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"log2sqs/config"
)

var q *sqs.SQS
var qURL = ""

// Buffered channel to trigger SQS reconnection
var sqsRestart = make(chan int, 1024)

// Open the SQS queue. Block until success because there is no point reading logs if
// there is nowhere to send them.
func openSQS() {
	for {
		err := connectSQS()
		if err != nil {
			log.Printf("Error opening queue: %s", err.Error())
			log.Printf("Sleeping for 30 seconds...")
			time.Sleep(30 * time.Second)
		} else {
			log.Printf("SQS queue %s opened", config.AWSQueueName)
			return
		}
	}
}

func connectSQS() error {
	var awsCredentials *credentials.Credentials
	var awsConfig *aws.Config

	// Initialize a session
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

	awsSession := session.Must(session.NewSession(awsConfig))
	q = sqs.New(awsSession)
	if q == nil {
		return errors.New("unable to create new AWS Session")
	}

	// Create a new request to list queues
	listQueuesRequest := sqs.ListQueuesInput{}
	listQueueResults, err := q.ListQueues(&listQueuesRequest)
	if err != nil {
		tmp := fmt.Sprintln("error listing SQS queues: ", err.Error())
		return errors.New(tmp)
	}

	// Search for requested queue name
	for _, t := range listQueueResults.QueueUrls {
		if strings.Contains(*t, config.AWSQueueName) {
			qURL = *t
			break
		}
	}

	if qURL == "" {
		tmp := fmt.Sprintf("unable to find SQS queue %s", config.AWSQueueName)
		return errors.New(tmp)
	}

	return nil
}

func sendSQS(msg []byte) error {

	//global.JSONPretty(msg) // For debugging only

	// Set up parameters
	var sendParams *sqs.SendMessageInput
	sendParams = &sqs.SendMessageInput{
		MessageBody: aws.String(string(msg)),
		QueueUrl:    aws.String(qURL),
	}

	// Send to SQS
	_, err := q.SendMessage(sendParams)
	if err != nil {
		// Request reconnection
		sqsRestart <- 1
		return err
	}

	return nil
}

// watchSQS waits for reconnection attempts and actions them
func watchSQS() {
	for {
		// Read message from channel
		// This is a blocking function
		_ = <-sqsRestart
		log.Printf("Received SQS reconnection request")

		// Open the SQS queue again
		openSQS()

		// Drain the channel to avoid unnecessary restarts
		for len(sqsRestart) > 0 {
			_ = <-sqsRestart
		}
	}
}
