package server

import "net"

type listener struct {
	conn net.Conn
}

func NewListener(conn net.Conn) *listener {
	return &listener{
		conn: conn,
	}
}
