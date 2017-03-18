package packet

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const Separator byte = '\n'

const innerSeparator byte = '|'
const innerSubseparator byte = ';'

const timeFormat = time.RFC3339Nano

type Data struct {
	PacketUUID  string      `json:"packet_uuid"`
	VehicleUUID string      `json:"vehicle_uuid"`
	Time        time.Time   `json:"time"`
	DroverData  DriverData  `json:"driver_data"`
	VehicleData VehicleData `json:"vehicle_data"`
}

func (p Data) Bytes() []byte {
	var bfr bytes.Buffer
	writeSeparatedValue(&bfr, []byte(p.PacketUUID))
	writeSeparatedValue(&bfr, []byte(p.VehicleUUID))
	writeSeparatedValue(&bfr, []byte(p.Time.Format(timeFormat)))
	writeSeparatedValue(&bfr, p.DroverData.Encode())
	writeSeparatedValue(&bfr, p.VehicleData.Encode())
	bfr.WriteByte(Separator)
	return bfr.Bytes()
}

func writeSeparatedValue(bfr *bytes.Buffer, val []byte) {
	if bfr.Len() != 0 {
		bfr.WriteByte(innerSeparator)
	}

	bfr.Write(val)
}

func writeSubseparatedValue(bfr *bytes.Buffer, val []byte) {
	if bfr.Len() != 0 {
		bfr.WriteByte(innerSubseparator)
	}

	bfr.Write(val)
}

func NewData(bts []byte) (*Data, error) {
	if len(bts) > 1 && bts[len(bts)-1] == Separator {
		bts = bts[:len(bts)-1]
	}

	pck := &Data{}

	parts := bytes.Split(bts, []byte{innerSeparator})
	const mustPartsCount = 5
	if len(parts) != mustPartsCount {
		return nil, fmt.Errorf("Data: Unexpected number of parts, expected %d, received %d.", mustPartsCount, len(parts))
	}

	var err error
	pck.PacketUUID = string(parts[0])
	pck.VehicleUUID = string(parts[1])
	pck.Time, err = time.Parse(timeFormat, string(parts[2]))
	if err != nil {
		return nil, err
	}

	err = pck.DroverData.Decode(parts[3])
	if err != nil {
		return nil, err
	}

	err = pck.VehicleData.Decode(parts[4])
	if err != nil {
		return nil, err
	}

	return pck, nil
}

func InitData() *Data {
	return &Data{
		PacketUUID: uuid.New().String(),
		Time:       time.Now(),
	}
}
