package network

import (
	"fmt"
	"net"
	"time"

	"github.com/skycoin/dmsg"
	"github.com/skycoin/dmsg/cipher"
	"github.com/skycoin/dmsg/noise"
)

const encryptHSTimout = 5 * time.Second

// Conn represents a network connection between two visors in skywire network
// This connection wraps raw network connection and is ready to use for sending data.
// It also provides skywire-specific methods on top of net.Conn
type Conn interface {
	net.Conn
	// LocalPK returns local public key of connection
	LocalPK() cipher.PubKey

	// RemotePK returns remote public key of connection
	RemotePK() cipher.PubKey

	// LocalPort returns local skywire port of connection
	// This is not underlying OS port, but port within skywire network
	LocalPort() uint16

	// RemotePort returns remote skywire port of connection
	// This is not underlying OS port, but port within skywire network
	RemotePort() uint16

	// Network returns network of connection
	Network() Type
}

type conn struct {
	net.Conn
	lAddr, rAddr dmsg.Addr
	freePort     func()
	connType     Type
}

func (c *conn) encrypt(lPK cipher.PubKey, lSK cipher.SecKey, initator bool) error {
	config := noise.Config{
		LocalPK:   lPK,
		LocalSK:   lSK,
		RemotePK:  c.rAddr.PK,
		Initiator: initator,
	}

	wrappedConn, err := EncryptConn(config, c.Conn)
	if err != nil {
		return fmt.Errorf("encrypt connection to %v@%v: %w", c.rAddr, c.Conn.RemoteAddr(), err)
	}

	c.Conn = wrappedConn
	return nil
}

// EncryptConn encrypts given connection
// todo: make private when tpconn.Conn is gone
func EncryptConn(config noise.Config, conn net.Conn) (net.Conn, error) {
	ns, err := noise.New(noise.HandshakeKK, config)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare stream noise object: %w", err)
	}

	wrappedConn, err := noise.WrapConn(conn, ns, encryptHSTimout)
	if err != nil {
		return nil, fmt.Errorf("error performing noise handshake: %w", err)
	}

	return wrappedConn, nil
}

// LocalAddr implements net.Conn
func (c *conn) LocalAddr() net.Addr {
	return c.lAddr
}

// RemoteAddr implements net.Conn
func (c *conn) RemoteAddr() net.Addr {
	return c.rAddr
}

// Close implements net.Conn
func (c *conn) Close() error {
	if c.freePort != nil {
		c.freePort()
	}

	return c.Conn.Close()
}

// LocalPK returns local public key of connection
func (c *conn) LocalPK() cipher.PubKey { return c.lAddr.PK }

// RemotePK returns remote public key of connection
func (c *conn) RemotePK() cipher.PubKey { return c.rAddr.PK }

// LocalPort returns local skywire port of connection
// This is not underlying OS port, but port within skywire network
func (c *conn) LocalPort() uint16 { return c.lAddr.Port }

// RemotePort returns remote skywire port of connection
// This is not underlying OS port, but port within skywire network
func (c *conn) RemotePort() uint16 { return c.rAddr.Port }

// Network returns network of connection
// todo: consider switching to Type instead of string
func (c *conn) Network() Type { return c.connType }
