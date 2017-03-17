package client

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type ConnManager struct {
	conns []*conn
	m     sync.Mutex
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
		ips = append(ips, c.ip)
	}

	return ips
}
