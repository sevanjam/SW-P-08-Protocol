package swp08

import (
	"fmt"
	"time"
)

// waitForAck waits for a 2-byte ACK or NAK: DLE 0x06 or DLE 0x15
// Returns true if ACK received, false if NAK or timeout/error
func waitForAck(state *ConnectionState) (bool, error) {
	if ACKWaitMode == NonBlocking {
		// Fire and forget: spawn a goroutine to listen for ACK
		go func() {
			select {
			case ack := <-state.ackChannel:
				if ack {
					fmt.Println("[SWP-08] ✅ ACK received (non-blocking)")
				} else {
					fmt.Println("[SWP-08] ❌ NAK received (non-blocking)")
				}
			case <-time.After(40 * time.Millisecond):
				fmt.Println("[SWP-08] ⚠ ACK/NAK timeout (non-blocking)")
			}
		}()
		// Return immediately
		return true, nil
	}

	// Blocking mode
	select {
	case ack := <-state.ackChannel:
		return ack, nil
	case <-time.After(2 * time.Second):
		return false, fmt.Errorf("timeout waiting for ACK/NAK")
	}
}
