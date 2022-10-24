package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/connorkuehl/netmem"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	listener, err := netmem.Listen()
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()

	go runEchoServer(ctx, listener)

	client, err := listener.(*netmem.Listener).Dial()
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer client.Close()

	// Print everything the server sends us in response to stdout.
	go func() {
		_, err = io.Copy(os.Stdout, client)
	}()

	// Send anything typed on stdin to the server.
	_, _ = io.Copy(client, os.Stdin)

	return nil
}

func runEchoServer(ctx context.Context, lis net.Listener) {
	conn, err := lis.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() && ctx.Err() == nil {
		fmt.Fprintln(conn, "Server says:", scanner.Text())
	}
}
