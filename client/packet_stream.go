package client

import (
	"encoding/binary"

	"io"
	"log"
)

type PacketStream struct {
	conn    io.ReadWriteCloser
	packets chan []byte
}

func NewPacketStream(conn io.ReadWriteCloser) *PacketStream {
	ps := PacketStream{
		conn:    conn,
		packets: make(chan []byte),
	}
	return &ps
}

func (p *PacketStream) readPackets() {
	go func() {
		for {
			var l uint32
			err := binary.Read(p.conn, binary.BigEndian, &l)
			if err != nil {
				log.Fatalln("Failed to read packet length:", err)
			}
			if l > 0 {
				packet := make([]byte, l)
				i, err := p.conn.Read(packet)
				if err != nil {
					log.Fatalln("Failed to read packet:", err)
				}
				if i != int(l) {
					log.Fatalln("Invalid packet size. Wanted:", l, "Read:", i)
				}
				p.packets <- packet
			}
		}
	}()
}

func (p *PacketStream) Read() []byte {
	return <-p.packets
}

func (p *PacketStream) Write(data []byte) (int, error) {
	err := binary.Write(p.conn, binary.BigEndian, uint32(len(data)))
	if err != nil {
		log.Fatalln("Failed to write packet length:", err)
	}
	return p.conn.Write(data)
}
