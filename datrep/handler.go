package datrep

import (
	"sync"

	"fmt"

	"github.com/2-guys-1-chick/c2c/network"
)

type handler struct {
	m sync.Mutex
}

func InitHandler() network.PacketHandler {
	return &handler{}
}

func (h *handler) Handle(packet *network.Packet) error {
	h.m.Lock()
	defer h.m.Unlock()
	// This may not be needed

	fmt.Printf("New incoming message: %s\n", packet.Text)
	return nil
}
