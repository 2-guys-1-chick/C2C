package client

import "net"

type conn struct {
	ip net.IP
	conn net.Conn
}