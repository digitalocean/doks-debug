package connections

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

type tlsConnector struct {
	targetHost string
}

// NewTLSConnector returns a TLS implementation of Connector
func NewTLSConnector(targetHost string) Connector {
	return &tlsConnector{targetHost: targetHost}
}

// Connect ...
func (c *tlsConnector) Connect(ctx context.Context) error {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", c.targetHost)
	if err != nil {
		return fmt.Errorf("net.ResolveTCPAddr failed: %s", err)
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return fmt.Errorf("net.DialTCP failed: %w", err)
	}
	defer tcpConn.Close()

	conn := tls.Client(tcpConn, conf)
	defer conn.Close()

	err = conn.Handshake()
	if err != nil {
		return fmt.Errorf("conn.Handshake failed: %w", err)
	}

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		return fmt.Errorf("conn.Write failed; numBytes: %d, error: %w", n, err)
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		return fmt.Errorf("conn.Read failed; numBytes: %d, error: %w", n, err)
	}

	return nil
}

// Type ...
func (c *tlsConnector) Type() string { return "tls" }
