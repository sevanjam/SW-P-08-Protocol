package swp08

import (
	"fmt"
	"time"
)

// waitForAck waits for a 2-byte ACK or NAK: DLE 0x06 or DLE 0x15
// Returns true if ACK received, false if NAK or timeout/error
func waitForAck(state *ConnectionState) (bool, error) {
	select {
	case ack := <-state.ackChannel:
		return ack, nil
	case <-time.After(2 * time.Second):
		return false, fmt.Errorf("timeout waiting for ACK/NAK")
	}
}
