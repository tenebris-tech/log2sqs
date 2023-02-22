//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package event

// Start initializes queues (internal and SQS) and starts the reading process
func Start() {

	// Initialize the queue
	initQueue()

	// Connect to SQS and block if required
	openSQS()

	// Start goroutines
	go watchSQS()
	go runQueue()
}
