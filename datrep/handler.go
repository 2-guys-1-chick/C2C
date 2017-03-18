package datrep

import (
	"sync"

	"encoding/json"

	"github.com/2-guys-1-chick/c2c/network"
	"github.com/2-guys-1-chick/c2c/network/packet"
	"github.com/2-guys-1-chick/c2c/network/ws"
)

type handler struct {
	m    sync.Mutex
	dist network.ByteDistributor
}

func InitHandler() network.PacketHandler {
	h := new(handler)
	wsSrv := ws.New()
	wsSrv.Run()
	h.dist = wsSrv
	return h
}

func (h *handler) Handle(packet *packet.Data) error {
	h.m.Lock()
	defer h.m.Unlock()

	jsonData, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	h.dist.Distribute(jsonData)
	return nil
}
