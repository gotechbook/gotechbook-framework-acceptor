package acceptor

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var forwardTables = map[string]struct {
	buf []byte
	err error
}{
	"test_handshake_type":     {[]byte{Handshake, 0x00, 0x00, 0x00}, nil},
	"test_handshake_ack_type": {[]byte{HandshakeAck, 0x00, 0x00, 0x00}, nil},
	"test_heartbeat_type":     {[]byte{Heartbeat, 0x00, 0x00, 0x00}, nil},
	"test_data_type":          {[]byte{Data, 0x00, 0x00, 0x00}, nil},
	"test_kick_type":          {[]byte{Kick, 0x00, 0x00, 0x00}, nil},
	"test_wrong_packet_type":  {[]byte{0x06, 0x00, 0x00, 0x00}, errors.New("wrong packet type")},
}

var (
	handshakeHeaderPacket = []byte{Handshake, 0x00, 0x00, 0x01, 0x01}
	invalidHeader         = []byte{0xff, 0x00, 0x00, 0x01}
)

var decodeTables = map[string]struct {
	data   []byte
	packet []*Packet
	err    error
}{
	"test_not_enough_bytes": {[]byte{0x01}, nil, nil},
	"test_error_on_forward": {invalidHeader, nil, errors.New("wrong packet type")},
	"test_forward":          {handshakeHeaderPacket, []*Packet{{Handshake, 1, []byte{0x01}}}, nil},
	"test_forward_many":     {append(handshakeHeaderPacket, handshakeHeaderPacket...), []*Packet{{Handshake, 1, []byte{0x01}}, {Handshake, 1, []byte{0x01}}}, nil},
}

func TestNewPacketCodec(t *testing.T) {
	t.Parallel()
	ppd := NewPacketCodec()
	assert.NotNil(t, ppd)
}

func TestForward(t *testing.T) {
	t.Parallel()
	for name, table := range forwardTables {
		t.Run(name, func(t *testing.T) {
			ppd := NewPacketCodec()
			sz, typ, err := ppd.forward(bytes.NewBuffer(table.buf))
			if table.err == nil {
				assert.Equal(t, Type(table.buf[0]), typ)
				assert.Equal(t, 0, sz)
			}
			assert.Equal(t, table.err, err)
		})
	}
}
func TestDecode(t *testing.T) {
	t.Parallel()

	for name, table := range decodeTables {
		t.Run(name, func(t *testing.T) {
			ppd := NewPacketCodec()
			packet, err := ppd.Decode(table.data)
			assert.Equal(t, table.err, err)
			assert.ElementsMatch(t, table.packet, packet)
		})
	}
}

func helperConcatBytes(packetType Type, length, data []byte) []byte {
	if data == nil {
		return nil
	}
	bytes := []byte{}
	bytes = append(bytes, byte(packetType))
	bytes = append(bytes, length...)
	bytes = append(bytes, data...)
	return bytes
}

var tooBigData = make([]byte, 1<<25)

var encodeTables = map[string]struct {
	packetType Type
	length     []byte
	data       []byte
	err        error
}{
	"test_encode_handshake":    {Handshake, []byte{0x00, 0x00, 0x02}, []byte{0x01, 0x00}, nil},
	"test_invalid_packet_type": {0xff, nil, nil, errors.New("wrong packet type")},
	"test_too_big_packet":      {Data, nil, tooBigData, errors.New("codec: packet size exceed")},
}

func TestEncode(t *testing.T) {
	t.Parallel()
	for name, table := range encodeTables {
		t.Run(name, func(t *testing.T) {
			ppe := NewPacketCodec()
			encoded, err := ppe.Encode(table.packetType, table.data)
			if table.err != nil {
				assert.Equal(t, table.err, err)
			} else {
				expectedEncoded := helperConcatBytes(table.packetType, table.length, table.data)
				assert.Equal(t, expectedEncoded, encoded)
			}
		})
	}
}
