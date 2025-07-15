package swp08

import (
	"fmt"
	"net"
	"time"
)

// waitForAck waits for a 2-byte ACK or NAK: DLE 0x06 or DLE 0x15
// Returns true if ACK received, false if NAK or timeout/error
func waitForAck(conn net.Conn) (bool, error) {
	buffer := make([]byte, 2)

	// Set a read deadline to avoid blocking forever
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetReadDeadline(time.Time{}) // Clear deadline after

	n, err := conn.Read(buffer)
	if err != nil {
		return false, fmt.Errorf("error reading ACK/NAK: %w", err)
	}
	if n != 2 {
		return false, fmt.Errorf("expected 2 bytes for ACK/NAK, got %d", n)
	}

	if buffer[0] != 0x10 {
		return false, fmt.Errorf("invalid first byte for ACK/NAK: 0x%02X", buffer[0])
	}

	switch buffer[1] {
	case 0x06:
		fmt.Println("[SWP-08] ✅ Received ACK from remote")
		return true, nil
	case 0x15:
		fmt.Println("[SWP-08] ❌ Received NAK from remote")
		return false, nil
	default:
		return false, fmt.Errorf("unknown ACK/NAK byte: 0x%02X", buffer[1])
	}
}
