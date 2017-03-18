package ws

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/2-guys-1-chick/c2c/cfg"
	"golang.org/x/net/websocket"
	"github.com/2-guys-1-chick/c2c/utils"
)

type srv struct {
	conns []*websocket.Conn
	m     sync.Mutex
}

func (s *srv) addConnection(conn *websocket.Conn) {
	s.m.Lock()
	defer s.m.Unlock()

	s.conns = append(s.conns, conn)
}

func (s *srv) removeConnection(conn1 *websocket.Conn) {
	s.m.Lock()
	defer s.m.Unlock()

	for i, conn2 := range s.conns {
		if conn1 == conn2 {
			s.conns = append(s.conns[:i], s.conns[i+1:]...)
			return
		}
	}

	// TODO better way
	fmt.Printf("Listener was not found: %v\n", errors.New("Not found"))
}

func New() *srv {
	s := new(srv)
	return s
}

func (s *srv) Run() {
	go s.listen()
}

func (s *srv) listen() error {
	http.Handle("/connect", websocket.Handler(s.handleNewConnection))
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetWsPort()), nil)
}

func (s *srv) handleNewConnection(ws *websocket.Conn) {
	s.addConnection(ws)

	for {
		var reply []byte
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			if utils.IsDisconnectError(err) {
				s.removeConnection(ws)
				ws.Close()
				break
			} else {
				fmt.Printf("WS Server: Unexpected server read error: %v\n", err)
			}
		}

		fmt.Println("INCOMIG MESSAGE", string(reply))
	}
}

func (s *srv) Distribute(msg []byte) {
	fmt.Println("LEN", len(s.conns))
	for _, conn := range s.conns {
		go s.sendMessage(msg, conn)
	}
}

func (s *srv) sendMessage(msg []byte, conn *websocket.Conn) {
	err := websocket.Message.Send(conn, string(msg))
	if err != nil {
		if utils.IsDisconnectError(err) {
			fmt.Printf("Client disconnected on read (%v)\n", err)
			s.removeConnection(conn)
			conn.Close()
		} else {
			fmt.Printf("Server: Unexpected read error: %v\n", err)
		}
	} else {
		fmt.Println("WS: Message sent successfully")
	}
}
