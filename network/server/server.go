package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"

	"log"

	"github.com/2-guys-1-chick/c2c/network"
	"github.com/2-guys-1-chick/c2c/network/packet"
	"github.com/2-guys-1-chick/c2c/utils"
)

// These client will receive updates
// Write to these
func StartServer(port int) (network.Distributor, error) {
	s, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	internalServer := &srv{
		srv: s,
	}

	go internalServer.acceptConnections()
	return internalServer, nil
}

func (s *srv) acceptConnections() {
	for {
		conn, err := s.srv.Accept()
		if err != nil {
			// handle error
			log.Println(err)
		}

		fmt.Println("Server: New incoming connection")
		l := &listener{
			conn: conn,
		}

		s.AddListener(l)
		go s.runConnectionCheck(l)
	}
}

type srv struct {
	srv       net.Listener
	listeners []*listener

	m sync.Mutex
}

func (s *srv) AddListener(l *listener) {
	s.m.Lock()
	defer s.m.Unlock()

	s.listeners = append(s.listeners, l)
}

func (s *srv) RemoveListener(l1 *listener) {
	s.m.Lock()
	defer s.m.Unlock()

	for i, l2 := range s.listeners {
		if l1 == l2 {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			return
		}
	}

	// TODO better way
	fmt.Printf("Listener was not found: %v\n", errors.New("Not found"))
}

func (s *srv) Distribute(packet *packet.Data) error {
	bts := packet.Bytes()

	for _, l := range s.listeners {
		_, err := l.conn.Write(bts)
		if err != nil {
			if utils.IsDisconnectError(err) {
				fmt.Printf("Client disconnected on write: %v\n", err)
				s.RemoveListener(l)
				l.conn.Close()
				break
			} else {
				fmt.Printf("Server: Unexpected write error: %v\n", err)
			}
		}
	}

	return nil
}

func (s *srv) runConnectionCheck(l *listener) {
	for {
		_, err := bufio.NewReader(l.conn).ReadString(packet.Separator)
		if err != nil {
			if utils.IsDisconnectError(err) {
				fmt.Printf("Client disconnected on read (%v)\n", err)
				s.RemoveListener(l)
				l.conn.Close()
				break
			} else {
				fmt.Printf("Server: Unexpected read error: %v\n", err)
			}
		}
	}
}

