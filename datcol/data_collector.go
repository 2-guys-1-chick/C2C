package datcol

import (
	"time"

	"fmt"

	"github.com/2-guys-1-chick/c2c/network"
)

type Collector struct {
	distributor network.Distributor
}

func (c *Collector) SetDistributor(distributor network.Distributor) {
	c.distributor = distributor
}

func (c *Collector) Run() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			packet := &network.Packet{}
			packet.Text = fmt.Sprintf("Hello, It is %s.", time.Now().Format("15:04:05"))
			c.distributor.Distribute(packet)
		}
	}
}
