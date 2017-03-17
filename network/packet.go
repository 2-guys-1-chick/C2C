package network

import "bytes"

const PacketSeparator byte = '\n'

type Packet struct {
}

func (p Packet) Bytes() ([]byte, error) {
	bfr := new(bytes.Buffer)
	// TODO Packet structure
	bfr.WriteByte(PacketSeparator)
	return bfr.Bytes(), nil
}

func NewPacket(bts []byte) (*Packet, error) {
	// TODO
	return &Packet{}, nil
}
