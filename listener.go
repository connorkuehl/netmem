package netmem

import (
	"fmt"
	"net"
	"sync"
)

// Listener is an in-memory net.Conn broker.
type Listener struct {
	toAccept chan net.Conn

	closeOnce sync.Once
	closed    chan struct{}
}

// Listen creates an in-memory net.Listener.
func Listen() (net.Listener, error) {
	l := &Listener{
		toAccept: make(chan net.Conn),
		closed:   make(chan struct{}),
	}

	return l, nil
}

func (l *Listener) isClosed() bool {
	select {
	case <-l.closed:
		return true
	default:
	}
	return false
}

// Dial allocates a new client connection.
func (l *Listener) Dial() (net.Conn, error) {
	client, server := net.Pipe()
	select {
	case l.toAccept <- server:
		return client, nil
	case <-l.closed:
		return nil, net.ErrClosed
	}
}

// Dialer returns a function that has the same shape as net.Dial.
func (l *Listener) Dialer() func(network, address string) (net.Conn, error) {
	return func(_, _ string) (net.Conn, error) {
		return l.Dial()
	}
}

// Accept blocks until an incoming connection is established.
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case c := <-l.toAccept:
		return c, nil
	case <-l.closed:
		return nil, net.ErrClosed
	}
}

// Close stops the Listener.
func (l *Listener) Close() error {
	err := net.ErrClosed
	l.closeOnce.Do(func() {
		close(l.closed)
		err = nil
	})
	return err
}

// Addr returns a net.Addr-compatible interface.
func (l *Listener) Addr() net.Addr {
	return Addr{
		vmaddr: fmt.Sprintf("%p", l),
	}
}
