package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"log"

	"syscall"

	"os"

	"github.com/2-guys-1-chick/c2c/network"
)

// These client will receive updates
// Write to these
func StartServer(port int) error {
	s, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	internalServer := &srv{
		srv: s,
	}

	go internalServer.acceptConnections()
	return nil
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

func (s *srv) WriteToAll(packet network.Packet) error {
	bts, err := packet.Bytes()
	if err != nil {
		return err
	}

	for _, l := range s.listeners {
		_, err := l.conn.Write(bts)
		if err != nil {
			if isDisconnectError(err) {
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
		_, err := bufio.NewReader(l.conn).ReadString(network.PacketSeparator)
		if err != nil {
			if isDisconnectError(err) {
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

func isDisconnectError(err error) bool {
	if err == io.EOF {
		return true
	}

	if netErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := netErr.Err.(*os.SyscallError); ok {
			if syscallErr.Err.Error() == syscall.ECONNRESET.Error() {
				return true
			}
		}
	}

	return false
}
