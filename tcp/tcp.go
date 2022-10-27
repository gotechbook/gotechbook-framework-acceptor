package tcp

import (
	"crypto/tls"
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	logger "github.com/gotechbook/gotechbook-framework-logger"
	"io"
	"io/ioutil"
	"net"
)

var _ acceptor.Acceptor = (*TCP)(nil)
var _ acceptor.Conn = (*tcpConn)(nil)

type TCP struct {
	addr     string
	connChan chan acceptor.Conn
	listener net.Listener
	running  bool
	certFile string
	keyFile  string
}

func NewTCP(addr string, certs ...string) *TCP {
	keyFile := ""
	certFile := ""
	if len(certs) != 2 && len(certs) != 0 {
		panic(acceptor.ErrInvalidCertificates)
	} else if len(certs) == 2 {
		certFile = certs[0]
		keyFile = certs[1]
	}

	return &TCP{
		addr:     addr,
		connChan: make(chan acceptor.Conn),
		running:  false,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (a *TCP) GetAddr() string {
	if a.listener != nil {
		return a.listener.Addr().String()
	}
	return ""
}
func (a *TCP) GetConnChan() chan acceptor.Conn {
	return a.connChan
}
func (a *TCP) Stop() {
	a.running = false
	a.listener.Close()
}
func (a *TCP) ListenAndServe() {
	if a.hasTLSCertificates() {
		a.ListenAndServeTLS(a.certFile, a.keyFile)
		return
	}
	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		logger.Log.Fatalf("Failed to listen: %s", err.Error())
	}
	a.listener = listener
	a.running = true
	a.serve()
}
func (a *TCP) ListenAndServeTLS(cert, key string) {
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logger.Log.Fatalf("Failed to listen: %s", err.Error())
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}

	listener, err := tls.Listen("tcp", a.addr, tlsCfg)
	if err != nil {
		logger.Log.Fatalf("Failed to listen: %s", err.Error())
	}
	a.listener = listener
	a.running = true
	a.serve()
}
func (a *TCP) hasTLSCertificates() bool {
	return a.certFile != "" && a.keyFile != ""
}
func (a *TCP) serve() {
	defer a.Stop()
	for a.running {
		conn, err := a.listener.Accept()
		if err != nil {
			logger.Log.Errorf("Failed to accept TCP connection: %s", err.Error())
			continue
		}
		a.connChan <- &tcpConn{
			Conn: conn,
		}
	}
}

type tcpConn struct {
	net.Conn
}

func (t *tcpConn) GetNextMessage() (b []byte, err error) {
	header, err := ioutil.ReadAll(io.LimitReader(t.Conn, acceptor.HeadLength))
	if err != nil {
		return nil, err
	}
	if len(header) == 0 {
		return nil, acceptor.ErrConnectionClosed
	}
	msgSize, _, err := acceptor.ParseHeader(header)
	if err != nil {
		return nil, err
	}
	msgData, err := ioutil.ReadAll(io.LimitReader(t.Conn, int64(msgSize)))
	if err != nil {
		return nil, err
	}
	if len(msgData) < msgSize {
		return nil, acceptor.ErrReceivedMsgSmallerThanExpected
	}
	return append(header, msgData...), nil
}
