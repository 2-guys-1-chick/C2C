package client

import (
	"bufio"
	"fmt"
	"net"

	"time"

	"github.com/2-guys-1-chick/c2c/network"
	"github.com/2-guys-1-chick/c2c/network/packet"
	"log"
)

func Connect(address string, port int, handler network.PacketHandler) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	go handleNewConnection(conn, handler)

	return conn, nil
}

func (cm *ConnManager) Connect(address string, port int) (net.Conn, error) {
	c, err := Connect(address, port, cm.handler)
	if err != nil {
		return nil, err
	}

	cm.createConnection(c)

	return c, nil
}

func (cm *ConnManager) SetPacketHandler(handler network.PacketHandler) {
	cm.handler = handler
}

func (cm *ConnManager) createConnection(c net.Conn) {
	cn := &conn{
		conn: c,
	}

	cm.addConnection(cn)
}

func (cm *ConnManager) RoundupConnect() error {
	return RoundupConnect(cm.getIPs(), cm.handler, cm.createConnection)

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

func handleNewConnection(conn net.Conn, handler network.PacketHandler) {
	defer conn.Close()
	for {
		pckBts, err := bufio.NewReader(conn).ReadBytes(packet.Separator)
		if err != nil {
			// handle error
		}

		go func(bts []byte) {
			err := handleBytes(pckBts, handler)
			if err != nil {
				log.Println(err)
				// handle error
			}
			// TODO do nothing?
		}(pckBts)
	}
}

func handleBytes(bts []byte, handler network.PacketHandler) error {
	pck, err := packet.NewData(bts)
	if err != nil {
		return err
	}

	return handler.Handle(pck)
}
