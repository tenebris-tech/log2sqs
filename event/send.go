//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package event

// Send is a public function to send directly via SQS without buffering
// This is useful for log files where buffering in memory doesn't make sense
func Send(msg []byte) error {
	return sendSQS(msg)
}
