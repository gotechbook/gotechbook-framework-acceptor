package acceptor

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	logger "github.com/gotechbook/gotechbook-framework-logger"
	"github.com/gotechbook/gotechbook-framework-net/constants"
	"net"
	"net/http"
)

var _ acceptor.Acceptor = (*WS)(nil)

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
		panic(constants.ErrInvalidCertificates)
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
func (w *WS) ListenAndServe() {
	if w.hasTLSCertificates() {
		w.ListenAndServeTLS(w.certFile, w.keyFile)
		return
	}

	var up = websocket.Upgrader{
		ReadBufferSize:  constants.IOBufferBytesSize,
		WriteBufferSize: constants.IOBufferBytesSize,
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
		ReadBufferSize:  constants.IOBufferBytesSize,
		WriteBufferSize: constants.IOBufferBytesSize,
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
