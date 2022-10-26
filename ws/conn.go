package acceptor

import (
	"github.com/gorilla/websocket"
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	"github.com/gotechbook/gotechbook-framework-net/constants"
	"github.com/gotechbook/gotechbook-framework-net/network/codec"
	"io"
	"net"
	"time"
)

var _ acceptor.Conn = (*WSConn)(nil)

type WSConn struct {
	conn   *websocket.Conn
	typ    int
	reader io.Reader
}

func NewWSConn(conn *websocket.Conn) (*WSConn, error) {
	c := &WSConn{conn: conn}
	return c, nil
}
func (c *WSConn) GetNextMessage() (b []byte, err error) {
	_, msgBytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if len(msgBytes) < constants.HeadLength {
		return nil, constants.ErrInvalidHeader
	}
	header := msgBytes[:constants.HeadLength]
	msgSize, _, err := codec.ParseHeader(header)
	if err != nil {
		return nil, err
	}
	dataLen := len(msgBytes[constants.HeadLength:])
	if dataLen < msgSize {
		return nil, constants.ErrReceivedMsgSmallerThanExpected
	} else if dataLen > msgSize {
		return nil, constants.ErrReceivedMsgBiggerThanExpected
	}
	return msgBytes, err
}
func (c *WSConn) Read(b []byte) (int, error) {
	if c.reader == nil {
		t, r, err := c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		c.typ = t
		c.reader = r
	}
	n, err := c.reader.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	} else if err == io.EOF {
		_, r, err := c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		c.reader = r
	}
	return n, nil
}
func (c *WSConn) Write(b []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}
func (c *WSConn) Close() error {
	return c.conn.Close()
}
func (c *WSConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *WSConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *WSConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}
func (c *WSConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *WSConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
