package acceptor

import (
	"errors"
	"fmt"
	"net"
)

type Type byte

const (
	_            Type = iota
	Handshake         = 0x01
	HandshakeAck      = 0x02
	Heartbeat         = 0x03
	Data              = 0x04
	Kick              = 0x05
)

const (
	HeadLength        = 4
	MaxPacketSize     = 1 << 24 //16MB
	IOBufferBytesSize = 4096
)

var (
	ErrInvalidCertificates            = errors.New("certificates must be exactly two")
	ErrWrongPacketType                = errors.New("wrong packet type")
	ErrPacketSizeExceed               = errors.New("codec: packet size exceed")
	ErrInvalidHeader                  = errors.New("invalid header")
	ErrReceivedMsgSmallerThanExpected = errors.New("received less data than expected, EOF")
	ErrReceivedMsgBiggerThanExpected  = errors.New("received more data than expected")
	ErrConnectionClosed               = errors.New("client connection closed")
)

type Acceptor interface {
	ListenAndServe()
	Stop()
	GetAddr() string
	GetConnChan() chan Conn
}

type Conn interface {
	GetNextMessage() (b []byte, err error)
	net.Conn
}

type Codec interface {
	Decode(data []byte) ([]*Packet, error)
	Encode(typType Type, data []byte) ([]byte, error)
}

type Packet struct {
	Type   Type
	Length int
	Data   []byte
}

func New() *Packet {
	return &Packet{}
}
func (p *Packet) String() string {
	return fmt.Sprintf("Type: %d, Length: %d, Data: %s", p.Type, p.Length, string(p.Data))
}
