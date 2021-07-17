//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func Send(msg string) error {

	// Set up parameters
	var sendParams *sqs.SendMessageInput
	sendParams = &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(qURL),
	}

	// Send to SQS
	_, err := q.SendMessage(sendParams)
	if err != nil {
		return err
	}

	return nil
}
