package acceptor

import (
	"github.com/gotechbook/gotechbook-framework-net/constants"
)

func ParseHeader(header []byte) (int, constants.Type, error) {
	if len(header) != constants.HeadLength {
		return 0, 0x00, constants.ErrInvalidHeader
	}
	typ := header[0]
	if typ < constants.Handshake || typ > constants.Kick {
		return 0, 0x00, constants.ErrWrongPacketType
	}
	size := BytesToInt(header[1:])
	if size > constants.MaxPacketSize {
		return 0, 0x00, constants.ErrPacketSizeExceed
	}
	return size, constants.Type(typ), nil
}

func BytesToInt(b []byte) int {
	result := 0
	for _, v := range b {
		result = result<<8 + int(v)
	}
	return result
}

func IntToBytes(n int) []byte {
	buf := make([]byte, 3)
	buf[0] = byte((n >> 16) & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte(n & 0xFF)
	return buf
}
