package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/2-guys-1-chick/c2c/network"
)

func Connect(address, port string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		// handle error
	}

	handleNewConnection(conn)
}

func handleNewConnection(conn net.Conn) {
	defer conn.Close()
	for {
		pckBts, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			// handle error
		}

		go func(bts []byte) {
			err := handleBytes(pckBts)
			if err != nil {
				// handle error
			}
			// TODO do nothing?
		}(pckBts)
	}
}

func handleBytes(bts []byte) error {
	packet, err := network.NewPacket(bts)
	if err != nil {
		return err
	}

	return handlePacket(packet)
}

func handlePacket(packet *network.Packet) error {
	// TODO handle packet
	return nil
}
