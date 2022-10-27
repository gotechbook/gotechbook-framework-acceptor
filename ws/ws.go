package ws

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	logger "github.com/gotechbook/gotechbook-framework-logger"
	"io"
	"net"
	"net/http"
	"time"
)

var _ acceptor.Acceptor = (*WS)(nil)
var _ acceptor.Conn = (*Conn)(nil)

type WS struct {
	addr     string
	connChan chan acceptor.Conn
	listener net.Listener
	certFile string
	keyFile  string
}

func NewWS(addr string, certs ...string) *WS {
	keyFile := ""
	certFile := ""
	if len(certs) != 2 && len(certs) != 0 {
		panic(acceptor.ErrInvalidCertificates)
	} else if len(certs) == 2 {
		certFile = certs[0]
		keyFile = certs[1]
	}
	w := &WS{
		addr:     addr,
		connChan: make(chan acceptor.Conn),
		certFile: certFile,
		keyFile:  keyFile,
	}
	return w
}

func (w *WS) ListenAndServe() {
	if w.hasTLSCertificates() {
		w.ListenAndServeTLS(w.certFile, w.keyFile)
		return
	}

	var up = websocket.Upgrader{
		ReadBufferSize:  acceptor.IOBufferBytesSize,
		WriteBufferSize: acceptor.IOBufferBytesSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	listener, err := net.Listen("tcp", w.addr)
	if err != nil {
		logger.Log.Fatalf("Failed to listen: %s", err.Error())
	}
	w.listener = listener

	w.serve(&up)
}
func (w *WS) Stop() {
	err := w.listener.Close()
	if err != nil {
		logger.Log.Errorf("Failed to stop: %s", err.Error())
	}
}
func (w *WS) GetAddr() string {
	if w.listener != nil {
		return w.listener.Addr().String()
	}
	return ""
}
func (w *WS) GetConnChan() chan acceptor.Conn {
	return w.connChan
}
func (w *WS) ListenAndServeTLS(cert, key string) {
	var up = websocket.Upgrader{
		ReadBufferSize:  acceptor.IOBufferBytesSize,
		WriteBufferSize: acceptor.IOBufferBytesSize,
	}

	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logger.Log.Fatalf("Failed to load x509: %s", err.Error())
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}
	listener, err := tls.Listen("tcp", w.addr, tlsCfg)
	if err != nil {
		logger.Log.Fatalf("Failed to listen: %s", err.Error())
	}
	w.listener = listener
	w.serve(&up)
}
func (w *WS) serve(up *websocket.Upgrader) {
	defer w.Stop()
	http.Serve(w.listener, &connHandler{
		up:       up,
		connChan: w.connChan,
	})
}
func (w *WS) hasTLSCertificates() bool {
	return w.certFile != "" && w.keyFile != ""
}

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

type connHandler struct {
	up       *websocket.Upgrader
	connChan chan acceptor.Conn
}

func (h *connHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	conn, err := h.up.Upgrade(rw, r, nil)
	if err != nil {
		logger.Log.Errorf("Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
		return
	}

	c, err := NewWSConn(conn)
	if err != nil {
		logger.Log.Errorf("Failed to create new ws connection: %s", err.Error())
		return
	}
	h.connChan <- c
}
