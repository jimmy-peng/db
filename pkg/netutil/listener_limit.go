package netutil

import (
	"errors"
	"net"
	"sync"
	"time"
)

// (etcd pkg.transport.limitListenerConn)
type connLimit struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

// ErrNotTCP defines reflection error for expected *net.TCPConn type.
//
// (etcd pkg.transport.ErrNotTCP)
var ErrNotTCP = errors.New("only tcp connections have keepalive")

func (c *connLimit) Close() error {
	err := c.Conn.Close()
	c.releaseOnce.Do(c.release)
	return err
}

func (c *connLimit) SetKeepAlive(doKeepAlive bool) error {
	tcpc, ok := c.Conn.(*net.TCPConn)
	if !ok {
		return ErrNotTCP
	}
	return tcpc.SetKeepAlive(doKeepAlive)
}

func (c *connLimit) SetKeepAlivePeriod(d time.Duration) error {
	tcpc, ok := c.Conn.(*net.TCPConn)
	if !ok {
		return ErrNotTCP
	}
	return tcpc.SetKeepAlivePeriod(d)
}

// (etcd pkg.transport.limitListener)
type listenerLimit struct {
	net.Listener
	sem chan struct{}
}

func (l *listenerLimit) acquire() {
	l.sem <- struct{}{}
}

func (l *listenerLimit) release() {
	<-l.sem
}

func (l *listenerLimit) Accept() (net.Conn, error) {
	l.acquire()
	c, err := l.Listener.Accept()
	if err != nil {
		l.release()
		return nil, err
	}
	return &connLimit{Conn: c, release: l.release}, nil
}

// NewListenerLimit returns a Listener that accepts at most n simultaneous
// connections from the provided Listener.
//
// (etcd pkg.transport.LimitListener)
func NewListenerLimit(l net.Listener, n int) net.Listener {
	return &listenerLimit{Listener: l, sem: make(chan struct{}, n)}
}
