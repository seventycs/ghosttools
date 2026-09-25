package utils

import (
	"net"
	"time"
)

func DialTimeout(host string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", host, timeout)
}