package swp08

import "net"

type ConnectionState struct {
	conn       net.Conn
	framer     Framer
	ackChannel chan bool
}
