package network

type Distributor interface {
	Distribute(packet *Packet) error
}

type PacketHandler interface {
	Handle(packet *Packet) error
}