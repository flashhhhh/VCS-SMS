package healthcheck

import (
	"net"
	"time"
)

func IsHostUp(address string) bool {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}