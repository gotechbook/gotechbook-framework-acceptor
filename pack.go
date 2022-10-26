package acceptor

import (
	"bytes"
	"errors"
)

type PacketCodec struct{}

func NewPacketCodec() *PacketCodec {
	return &PacketCodec{}
}

func (c *PacketCodec) forward(buf *bytes.Buffer) (int, Type, error) {
	header := buf.Next(HeadLength)
	return ParseHeader(header)
}
func (c *PacketCodec) Decode(data []byte) ([]*Packet, error) {
	buf := bytes.NewBuffer(nil)
	buf.Write(data)

	var (
		packets []*Packet
		err     error
	)
	// check length
	if buf.Len() < HeadLength {
		return nil, nil
	}

	// first time
	size, typ, err := c.forward(buf)
	if err != nil {
		return nil, err
	}

	for size <= buf.Len() {
		p := &Packet{Type: typ, Length: size, Data: buf.Next(size)}
		packets = append(packets, p)

		// if no more packets, break
		if buf.Len() < HeadLength {
			break
		}

		size, typ, err = c.forward(buf)
		if err != nil {
			return nil, err
		}
	}

	return packets, nil
}
func (c *PacketCodec) Encode(typ Type, data []byte) ([]byte, error) {
	if typ < Handshake || typ > Kick {
		return nil, errors.New("wrong packet type")
	}
	if len(data) > MaxPacketSize {
		return nil, errors.New("codec: packet size exceed")
	}
	p := &Packet{Type: typ, Length: len(data)}
	buf := make([]byte, p.Length+HeadLength)
	buf[0] = byte(p.Type)

	copy(buf[1:HeadLength], IntToBytes(p.Length))
	copy(buf[HeadLength:], data)

	return buf, nil
}
