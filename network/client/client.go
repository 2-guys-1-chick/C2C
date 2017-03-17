package client

import (
	"bufio"
	"fmt"
	"net"

	"time"

	"github.com/2-guys-1-chick/c2c/network"
)

func Connect(address string, port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	go handleNewConnection(conn)

	return conn, nil
}

func (cm *ConnManager) Connect(address string, port int) (net.Conn, error) {
	c, err := Connect(address, port)
	if err != nil {
		return nil, err
	}

	cm.createConnection(c)

	return c, nil
}

func (cm *ConnManager) createConnection(c net.Conn) {
	cn := &conn{
		conn: c,
	}

	cm.addConnection(cn)
}

func (cm *ConnManager) RoundupConnect() error {
	return RoundupConnect(cm.getIPs(), cm.createConnection)

}

func (cm *ConnManager) InitRoundup() {
	ticker := time.NewTicker(30 * time.Second) // TODO constant
	cm.RoundupConnect()
	go func() {
		for {
			select {
			case <-ticker.C:
				cm.RoundupConnect()
			}
		}
	}()
}

func handleNewConnection(conn net.Conn) {
	defer conn.Close()
	for {
		pckBts, err := bufio.NewReader(conn).ReadBytes(network.PacketSeparator)
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
