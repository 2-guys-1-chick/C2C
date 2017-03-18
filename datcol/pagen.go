package datcol

import (
	"time"

	"github.com/2-guys-1-chick/c2c/network/packet"
	"github.com/google/uuid"
)

type packetGenerator struct {
	vehicleUUID   string
	initializedAt *time.Time
}

func (g *packetGenerator) init() {
	now := time.Now()
	g.initializedAt = &now

	g.vehicleUUID = uuid.New().String()

}
func (g *packetGenerator) GetNext() *packet.Data {
	if g.initializedAt == nil {
		g.init()
	}

	data := packet.InitData()
	data.VehicleUUID = g.vehicleUUID

	diff := time.Now().Sub(*g.initializedAt)
	pnt, speed := calculateMovement(int(diff.Seconds() * 1000))
	data.VehicleData.Speed = speed
	data.VehicleData.Geo = *pnt
	return data
}
