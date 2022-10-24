package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/connorkuehl/netmem"
)

func main() {
	l, _ := netmem.Listen()
	defer l.Close()

	go func() {
		conn, _ := l.Accept()
		defer conn.Close()
		_, _ = io.WriteString(conn, "Greetings, friend!")
	}()

	cli, _ := l.(*netmem.Listener).Dial()
	defer cli.Close()

	var rsp strings.Builder
	_, _ = io.Copy(&rsp, cli)
	fmt.Printf("Incoming transmission: %q\n", rsp.String())
}
