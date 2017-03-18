package packet_test

import (
	"testing"

	"time"

	"fmt"

	"github.com/2-guys-1-chick/c2c/network/packet"
	"github.com/kellydunn/golang-geo"
	. "gopkg.in/check.v1"
)

type packetSuite struct {
	testPairs []packetTestPair
}

var _ = Suite(new(packetSuite))

func TestPacket(t *testing.T) { TestingT(t) }

func (s *packetSuite) SetUpSuite(c *C) {
	pck := packet.Data{
		PacketUUID:  "packet-uuid",
		VehicleUUID: "vehicle-uuid",
		Time:        time.Date(2017, 1, 1, 10, 0, 0, 820366527, time.Local),
	}

	s.testPairs = append(s.testPairs, packetTestPair{
		p: pck,
		s: "packet-uuid|vehicle-uuid|2017-01-01T10:00:00.820366527+01:00||0.000;0.0000000,0.0000000;0.0000;0.0000;",
	})

	pck.VehicleData = packet.VehicleData{
		Speed:     60,
		Geo:       *geo.NewPoint(12.5432, -1.123),
		Weight:    1.02,
		TireWear:  0.80,
		DriveMode: packet.DriveModeAutopilot,
	}

	s.testPairs = append(s.testPairs, packetTestPair{
		p: pck,
		s: "packet-uuid|vehicle-uuid|2017-01-01T10:00:00.820366527+01:00||60.000;12.5432000,-1.1230000;1.0200;0.8000;AUTO",
	})

	pck.VehicleData = packet.VehicleData{}
	pck.DroverData = packet.DriverData{
		Moods: packet.Moods(packet.MoodTired, packet.MoodImpetuous),
	}

	s.testPairs = append(s.testPairs, packetTestPair{
		p: pck,
		s: "packet-uuid|vehicle-uuid|2017-01-01T10:00:00.820366527+01:00|TIRED,IMPETUOUS|0.000;0.0000000,0.0000000;0.0000;0.0000;",
	})

	s.testPairs = append(s.testPairs, packetTestPair{
		p: pck,
		s: "packet-uuid|vehicle-uuid|2017-01-01T10:00:00.820366527+01:00|TIRED,IMPETUOUS,INVALID|0.000;0.0000000,0.0000000;0.0000;0.0000;",
	})
}

type packetTestPair struct {
	p packet.Data
	s string
}

func (s *packetSuite) TestEncode(c *C) {
	for _, pair := range s.testPairs {
		res := pair.p.Bytes()
		fmt.Println(string(res))
		res = res[:len(res)-1]
		c.Assert(res, DeepEquals, []byte(pair.s))
	}
}

func (s *packetSuite) TestDecode(c *C) {
	for _, pair := range s.testPairs {
		data, err := packet.NewData([]byte(pair.s + string(packet.Separator)))
		c.Assert(err, IsNil)
		c.Assert(*data, DeepEquals, pair.p)
	}
}
