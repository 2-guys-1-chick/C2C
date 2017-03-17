package client

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/2-guys-1-chick/c2c/network"
)

type ConnManager struct {
	conns   []*conn
	m       sync.Mutex
	handler network.PacketHandler
}

func (cm *ConnManager) addConnection(conn *conn) {
	cm.m.Lock()
	defer cm.m.Unlock()

	cm.conns = append(cm.conns, conn)
}

func (cm *ConnManager) removeConnection(conn1 *conn) {
	cm.m.Lock()
	defer cm.m.Unlock()

	for i, conn2 := range cm.conns {
		if conn1 == conn2 {
			cm.conns = append(cm.conns[:i], cm.conns[i+1:]...)
			return
		}
	}

	// TODO better way
	fmt.Printf("Listener was not found: %v\n", errors.New("Not found"))
}

func (cm *ConnManager) getIPs() []net.IP {
	var ips []net.IP
	for _, c := range cm.conns {
		if addr, ok := c.conn.RemoteAddr().(*net.TCPAddr); ok {
			ips = append(ips, addr.IP)
		}
	}

	return ips
}
