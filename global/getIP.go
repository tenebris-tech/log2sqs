//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package global

import (
	"fmt"
	"net"
)

// GetOutboundIP uses net.Dial to determine the hosts preferred IP address
// This makes it easy to avoid returning localhost, etc.
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	return fmt.Sprint(conn.LocalAddr().(*net.UDPAddr).IP)
}
