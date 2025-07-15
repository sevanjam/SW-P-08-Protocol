package swp08

import (
	"fmt"
)

// Entry point for handling a valid SW-P-08 message
func handleMessage(state *ConnectionState, msg []byte) {
	if len(msg) < 2 {
		fmt.Println("[SWP-08] ❌ Message too short to contain command and byte count")
		return
	}

	command := msg[0]
	// byteCount := data[len(data)-2]
	message := msg[1 : len(msg)-2] // Exclude command, byte count, checksum

	fmt.Printf("[SWP-08] ➡ Handling command: 0x%02X\n", command)
	fmt.Printf("[SWP-08] ↪ Message body (%d bytes):\n", len(message))
	printHex(message)

	// Future: check byteCount again here if needed (already validated earlier)

	switch command {

	//
	case 0x02:
		handleSetCrosspoint(state, message)

	// Dual Controller Status Request Message
	case 0x08:
		handleDualControllerStatus(state)

	case 0x15:
		fmt.Println("[SWP-08] 📥 Crosspoint Tally Dump Request (command 0x15)")
		if len(message) < 1 {
			fmt.Println("[SWP-08] ❌ Invalid message length for Tally Dump Request")
			return
		}
		handleTallyDumpRequest(state, message[0])

	default:
		fmt.Printf("[SWP-08] ❓ Unknown or unhandled command: 0x%02X\n", command)
	}
}
