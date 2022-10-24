package netmem

import (
	"errors"
	"io"
	"net"
	"strings"
	"testing"
)

func TestAccept(t *testing.T) {
	listener, err := Listen()
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()

	go func() {
		cli, err := listener.(*Listener).Dial()
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(cli, strings.NewReader("hi"))
		if err != nil {
			t.Error(err)
		}
		cli.Close()
	}()

	c, err := listener.Accept()
	if err != nil {
		t.Error(err)
	}

	var got strings.Builder
	_, err = io.Copy(&got, c)
	if err != nil {
		t.Error(err)
	}

	if got.String() != "hi" {
		t.Errorf("wanted %q, got %q", "hi", got.String())
	}
}

func TestClose(t *testing.T) {
	listener, err := Listen()
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()

	lis := listener.(*Listener)
	if lis.isClosed() {
		t.Errorf("listener is unexpectedly closed")
	}

	err = listener.Close()
	if err != nil {
		t.Error(err)
	}

	if !lis.isClosed() {
		t.Errorf("listener is unexpectedly open")
	}

	err = listener.Close()
	if !errors.Is(err, net.ErrClosed) {
		t.Errorf("want %v, got %v", net.ErrClosed, err)
	}

	_, err = listener.Accept()
	if !errors.Is(err, net.ErrClosed) {
		t.Errorf("want %v, got %v", net.ErrClosed, err)
	}
}

func TestAddr(t *testing.T) {
	listener, err := Listen()
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()

	a := listener.Addr()
	if a.Network() != "memory" {
		t.Errorf("want %q, got %q", "memory", a.Network())
	}
}
