package tcp

import (
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	utils "github.com/gotechbook/gotechbook-framework-utils"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var tcpAcceptorTables = []struct {
	name     string
	addr     string
	certs    []string
	panicErr error
}{
	{"test_1", "0.0.0.0:0", []string{"../fixtures/server.crt", "../fixtures/server.key"}, nil},
	{"test_2", "0.0.0.0:0", []string{}, nil},
	{"test_3", "127.0.0.1:0", []string{"wqd"}, acceptor.ErrInvalidCertificates},
	{"test_4", "127.0.0.1:0", []string{"wqd", "wqdqwd", "wqdqdqwd"}, acceptor.ErrInvalidCertificates},
}

func TestNewTCPGetConnChanAndGetAddr(t *testing.T) {
	t.Parallel()
	for _, table := range tcpAcceptorTables {
		t.Run(table.name, func(t *testing.T) {
			if table.panicErr != nil {
				assert.PanicsWithValue(t, table.panicErr, func() {
					NewTCP(table.addr, table.certs...)
				})
			} else {
				var a *TCP
				assert.NotPanics(t, func() {
					a = NewTCP(table.addr, table.certs...)
				})

				if len(table.certs) == 2 {
					assert.Equal(t, table.certs[0], a.certFile)
					assert.Equal(t, table.certs[1], a.keyFile)
				} else {
					assert.Equal(t, "", a.certFile)
					assert.Equal(t, "", a.keyFile)
				}
				assert.NotNil(t, a)
			}
		})
	}
}

func TestGetAddr(t *testing.T) {
	t.Parallel()
	for _, table := range tcpAcceptorTables {
		t.Run(table.name, func(t *testing.T) {
			a := NewTCP(table.addr)
			// returns nothing because not listening yet
			assert.Equal(t, "", a.GetAddr())
		})
	}
}

func TestGetConnChan(t *testing.T) {
	t.Parallel()
	for _, table := range tcpAcceptorTables {
		t.Run(table.name, func(t *testing.T) {
			a := NewTCP(table.addr)
			assert.NotNil(t, a.GetConnChan())
		})
	}
}

func TestListenAndServe(t *testing.T) {
	for _, table := range tcpAcceptorTables {
		t.Run(table.name, func(t *testing.T) {
			a := NewTCP(table.addr)
			defer a.Stop()
			c := a.GetConnChan()
			go a.ListenAndServe()
			// should be able to connect within 100 milliseconds
			utils.ShouldEventuallyReturn(t, func() error {
				n, err := net.Dial("tcp", a.GetAddr())
				defer n.Close()
				return err
			}, nil, 10*time.Millisecond, 100*time.Millisecond)
			conn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond)
			assert.NotNil(t, conn)
		})
	}
}

func TestListenAndServeTLS(t *testing.T) {
	for _, table := range tcpAcceptorTables {
		t.Run(table.name, func(t *testing.T) {
			a := NewTCP(table.addr)
			defer a.Stop()
			c := a.GetConnChan()

			go a.ListenAndServeTLS("../fixtures/server.crt", "../fixtures/server.key")
			// should be able to connect within 100 milliseconds
			utils.ShouldEventuallyReturn(t, func() error {
				n, err := net.Dial("tcp", a.GetAddr())
				defer n.Close()
				return err
			}, nil, 10*time.Millisecond, 100*time.Millisecond)
			conn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond)
			assert.NotNil(t, conn)
		})
	}
}

func TestStop(t *testing.T) {
	for _, table := range tcpAcceptorTables {
		t.Run(table.name, func(t *testing.T) {
			a := NewTCP(table.addr)
			go a.ListenAndServe()
			// should be able to connect within 100 milliseconds
			utils.ShouldEventuallyReturn(t, func() error {
				_, err := net.Dial("tcp", a.GetAddr())
				return err
			}, nil, 10*time.Millisecond, 100*time.Millisecond)
			a.Stop()
			_, err := net.Dial("tcp", table.addr)
			assert.Error(t, err)
		})
	}
}

func TestGetNextMessage(t *testing.T) {
	tables := []struct {
		name string
		data []byte
		err  error
	}{
		{"invalid_header", []byte{0x00, 0x00, 0x00, 0x00}, acceptor.ErrWrongPacketType},
		{"valid_message", []byte{0x02, 0x00, 0x00, 0x01, 0x00}, nil},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			a := NewTCP("0.0.0.0:0")
			go a.ListenAndServe()
			defer a.Stop()
			c := a.GetConnChan()
			// should be able to connect within 100 milliseconds
			var conn net.Conn
			var err error
			utils.ShouldEventuallyReturn(t, func() error {
				conn, err = net.Dial("tcp", a.GetAddr())
				return err
			}, nil, 10*time.Millisecond, 100*time.Millisecond)

			defer conn.Close()
			playerConn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond).(acceptor.Conn)
			_, err = conn.Write(table.data)
			assert.NoError(t, err)

			msg, err := playerConn.GetNextMessage()
			if table.err != nil {
				assert.EqualError(t, err, table.err.Error())
			} else {
				assert.Equal(t, table.data, msg)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetNextMessageTwoMessagesInBuffer(t *testing.T) {
	a := NewTCP("0.0.0.0:0")
	go a.ListenAndServe()
	defer a.Stop()
	c := a.GetConnChan()
	// should be able to connect within 100 milliseconds
	var conn net.Conn
	var err error
	utils.ShouldEventuallyReturn(t, func() error {
		conn, err = net.Dial("tcp", a.GetAddr())
		return err
	}, nil, 10*time.Millisecond, 100*time.Millisecond)
	defer conn.Close()

	playerConn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond).(acceptor.Conn)
	msg1 := []byte{0x01, 0x00, 0x00, 0x01, 0x02}
	msg2 := []byte{0x02, 0x00, 0x00, 0x02, 0x01, 0x01}
	buffer := append(msg1, msg2...)
	_, err = conn.Write(buffer)
	assert.NoError(t, err)

	msg, err := playerConn.GetNextMessage()
	assert.NoError(t, err)
	assert.Equal(t, msg1, msg)

	msg, err = playerConn.GetNextMessage()
	assert.NoError(t, err)
	assert.Equal(t, msg2, msg)
}

func TestGetNextMessageEOF(t *testing.T) {
	a := NewTCP("0.0.0.0:0")
	go a.ListenAndServe()
	defer a.Stop()
	c := a.GetConnChan()
	// should be able to connect within 100 milliseconds
	var conn net.Conn
	var err error
	utils.ShouldEventuallyReturn(t, func() error {
		conn, err = net.Dial("tcp", a.GetAddr())
		return err
	}, nil, 10*time.Millisecond, 100*time.Millisecond)

	playerConn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond).(acceptor.Conn)
	buffer := []byte{0x02, 0x00, 0x00, 0x02, 0x01}
	_, err = conn.Write(buffer)
	assert.NoError(t, err)

	go func() {
		time.Sleep(100 * time.Millisecond)
		conn.Close()
	}()

	_, err = playerConn.GetNextMessage()
	assert.EqualError(t, err, acceptor.ErrReceivedMsgSmallerThanExpected.Error())
}

func TestGetNextMessageEmptyEOF(t *testing.T) {
	a := NewTCP("0.0.0.0:0")
	go a.ListenAndServe()
	defer a.Stop()
	c := a.GetConnChan()
	// should be able to connect within 100 milliseconds
	var conn net.Conn
	var err error
	utils.ShouldEventuallyReturn(t, func() error {
		conn, err = net.Dial("tcp", a.GetAddr())
		return err
	}, nil, 10*time.Millisecond, 100*time.Millisecond)

	playerConn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond).(acceptor.Conn)

	go func() {
		time.Sleep(100 * time.Millisecond)
		conn.Close()
	}()

	_, err = playerConn.GetNextMessage()
	assert.EqualError(t, err, acceptor.ErrConnectionClosed.Error())
}

func TestGetNextMessageInParts(t *testing.T) {
	a := NewTCP("0.0.0.0:0")
	go a.ListenAndServe()
	defer a.Stop()
	c := a.GetConnChan()
	// should be able to connect within 100 milliseconds
	var conn net.Conn
	var err error
	utils.ShouldEventuallyReturn(t, func() error {
		conn, err = net.Dial("tcp", a.GetAddr())
		return err
	}, nil, 10*time.Millisecond, 100*time.Millisecond)

	defer conn.Close()
	playerConn := utils.ShouldEventuallyReceive(t, c, 100*time.Millisecond).(acceptor.Conn)
	part1 := []byte{0x02, 0x00, 0x00, 0x03, 0x01}
	part2 := []byte{0x01, 0x02}
	_, err = conn.Write(part1)
	assert.NoError(t, err)

	go func() {
		time.Sleep(200 * time.Millisecond)
		_, err = conn.Write(part2)
	}()

	msg, err := playerConn.GetNextMessage()
	assert.NoError(t, err)
	assert.Equal(t, msg, append(part1, part2...))

}
