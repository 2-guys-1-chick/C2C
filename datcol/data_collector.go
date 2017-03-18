package datcol

import (
	"time"

	"fmt"

	"github.com/2-guys-1-chick/c2c/network"
	"github.com/2-guys-1-chick/c2c/network/packet"
	"github.com/kellydunn/golang-geo"
	"github.com/google/uuid"
)

type Collector struct {
	distributor network.Distributor
	generator   packetGenerator
}

func (c *Collector) SetDistributor(distributor network.Distributor) {
	c.distributor = distributor
}

func (c *Collector) Run() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			data := c.generator.GetNext()
			if data == nil {
				ticker.Stop()
			}
			fmt.Print("Sending: ", string(data.Bytes()))
			c.distributor.Distribute(data)
		}
	}
}

type packetGenerator struct {
	vehicleUUID string
	initialized bool
}

func (g *packetGenerator) init() {
	defer func() {
		g.initialized = true
	}()

	g.vehicleUUID = uuid.New().String()

}
func (g *packetGenerator) GetNext() *packet.Data {
	if !g.initialized {
		g.init()
	}

	data := packet.InitData()
	data.VehicleUUID = g.vehicleUUID
	data.VehicleData.Speed = 50
	data.VehicleData.Geo = *geo.NewPoint(12.312312, -1.1231)
	return data
}
