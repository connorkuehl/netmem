package netmem

// Addr is a net.Addr-compliant struct.
type Addr struct {
	vmaddr string
}

// Network is always "memory"
func (a Addr) Network() string { return "memory" }

// String is the hexadecimal-formatted string representation of the
// Listener's virtual memory address.
func (a Addr) String() string { return a.vmaddr }
