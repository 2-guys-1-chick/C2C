package datcol

import (
	"time"

	"github.com/2-guys-1-chick/c2c/network"
)

type Collector struct {
	distributor network.Distributor
	generator   packetGenerator
}

func (c *Collector) SetDistributor(distributor network.Distributor) {
	c.distributor = distributor
}

func (c *Collector) Run(done chan<- struct{}) {
	ticker := time.NewTicker(200 * time.Millisecond)

Loop:
	for {
		select {
		case <-ticker.C:
			data := c.generator.GetNext()
			if data == nil {
				ticker.Stop()
				break Loop

			}
			//fmt.Print("Sending: ", string(data.Bytes()))
			c.distributor.Distribute(data)
		}
	}

	done <- struct{}{}
}
