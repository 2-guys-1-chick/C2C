package datcol

import (
	"time"

	"github.com/2-guys-1-chick/c2c/cfg"
	"github.com/2-guys-1-chick/c2c/network/packet"
)

type packetGenerator struct {
	vehicleUUID   string
	initializedAt *time.Time
}

func (g *packetGenerator) init() {
	now := time.Now()
	g.initializedAt = &now

	g.vehicleUUID = cfg.GetVehicleId()

}
func (g *packetGenerator) GetNext() *packet.Data {
	if g.initializedAt == nil {
		g.init()
	}

	data := packet.InitData()
	data.VehicleUUID = g.vehicleUUID

	diff := time.Now().Sub(*g.initializedAt)
	pnt, speed := calculateMovement(g.vehicleUUID, int(diff.Seconds()*1000))
	if pnt == nil {
		return nil
	}

	data.VehicleData.Speed = speed
	data.VehicleData.Geo = *pnt
	return data
}
