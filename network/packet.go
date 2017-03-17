package network

import "bytes"

const PacketSeparator byte = '\n'

type Packet struct {
	Text string
}

func (p Packet) Bytes() ([]byte, error) {
	bfr := new(bytes.Buffer)
	// TODO Packet structure
	bfr.WriteString(p.Text)
	bfr.WriteByte(PacketSeparator)
	return bfr.Bytes(), nil
}

func NewPacket(bts []byte) (*Packet, error) {
	// TODO
	pck := &Packet{}
	pck.Text = string(bts)
	return pck, nil
}
