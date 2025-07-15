package swp08

import (
	"fmt"
	"net"
)

// Send DLE ACK (0x10 0x06)
func sendAcknowledge(conn net.Conn) {
	ack := []byte{0x10, 0x06}
	_, err := conn.Write(ack)
	if err != nil {
		fmt.Println("[SWP-08] ❌ Failed to send ACK:", err)
	} else {
		fmt.Println("[SWP-08] ✅ Sent ACK")
	}
}

// Send DLE NAK (0x10 0x15)
func sendNegativeAcknowledge(conn net.Conn) {
	nak := []byte{0x10, 0x15}
	_, err := conn.Write(nak)
	if err != nil {
		fmt.Println("[SWP-08] ❌ Failed to send NAK:", err)
	} else {
		fmt.Println("[SWP-08] ❌ Sent NAK")
	}
}
