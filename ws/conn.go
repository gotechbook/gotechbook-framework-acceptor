package ws

import (
	"github.com/gorilla/websocket"
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	"io"
	"net"
	"time"
)

var _ acceptor.Conn = (*Conn)(nil)

type Conn struct {
	conn   *websocket.Conn
	typ    int
	reader io.Reader
}

func NewWSConn(conn *websocket.Conn) (*Conn, error) {
	c := &Conn{conn: conn}
	return c, nil
}
func (c *Conn) GetNextMessage() (b []byte, err error) {
	_, msgBytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if len(msgBytes) < acceptor.HeadLength {
		return nil, acceptor.ErrInvalidHeader
	}
	header := msgBytes[:acceptor.HeadLength]
	msgSize, _, err := acceptor.ParseHeader(header)
	if err != nil {
		return nil, err
	}
	dataLen := len(msgBytes[acceptor.HeadLength:])
	if dataLen < msgSize {
		return nil, acceptor.ErrReceivedMsgSmallerThanExpected
	} else if dataLen > msgSize {
		return nil, acceptor.ErrReceivedMsgBiggerThanExpected
	}
	return msgBytes, err
}
func (c *Conn) Read(b []byte) (int, error) {
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
func (c *Conn) Write(b []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}
func (c *Conn) Close() error {
	return c.conn.Close()
}
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *Conn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
