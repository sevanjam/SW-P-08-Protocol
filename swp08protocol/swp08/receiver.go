package swp08

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	state := &ConnectionState{
		conn:       conn,
		framer:     Framer{},
		ackChannel: make(chan bool, 1), // buffered channel to avoid blocking
	}

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// handle error or connection close here
			break
		}

		data := buffer[:n]

		fmt.Printf("[SWP-08] 📥 Received %d bytes:\n", n)
		printHex(data)

		// Check for ACK/NAK first
		if len(data) == 2 && data[0] == 0x10 {
			switch data[1] {
			case 0x06:
				fmt.Println("[SWP-08] ✅ Received ACK")
				state.ackChannel <- true
				continue
			case 0x15:
				fmt.Println("[SWP-08] ❌ Received NAK")
				state.ackChannel <- false
				continue
			}
		}

		messages := state.framer.Feed(data)
		for _, msg := range messages {
			handleMessage(state, msg) // pass *ConnectionState, not net.Conn
		}
	}
}
