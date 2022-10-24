# netmem

Package netmem provides a completely in-memory Listener that brokers net.Conns
based on net.Pipe.

Its primary use case is to enable broader end-to-end testing functionality
from within the scope of a unit test. For example, to test the public
implementation of a client package against a fully-functioning instance of
the server package.

Generally speaking, if your server implementation allows you to inject
its net.Listener, you can inject this instead and it will accept in-memory
network connections.

```go
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
```

```console
$ go run ./examples/readme
Incoming transmission: "Greetings, friend!"
```
